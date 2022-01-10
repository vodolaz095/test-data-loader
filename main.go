package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// DataElement is minimal element of data
type DataElement struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// OutputDataStructure is data structure format output to/from json files
type OutputDataStructure struct {
	Data []DataElement `json:"data"`
}

// Concurrency depicts number of concurrent processes being executed
const Concurrency = 12

// Output is used for output
var Output map[int]DataElement

var channelForFilesToParse chan string
var channelForDataElementsToSave chan DataElement

var inputDirectory string
var outputFileName string
var ignoreDuplicates bool

// ReadDirectory reads directory
func ReadDirectory(pathToDirectory string) (err error) {
	err = filepath.Walk(pathToDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Source directory not found")
				os.Exit(10)
			}
			return err
		}

		if filepath.Dir(path) == pathToDirectory {
			if !info.IsDir() && filepath.Ext(info.Name()) == ".json" {
				channelForFilesToParse <- path
			}
		}
		return nil
	})
	return
}

// Parse parses
func Parse(pathToFile string) (err error) {
	var elements OutputDataStructure
	data, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &elements)
	if err != nil {
		return
	}
	for _, el := range elements.Data {
		channelForDataElementsToSave <- el
	}
	return
}

func main() {
	var err error
	flag.StringVar(&inputDirectory, "source-dir", "", "The path to the directory to read files from")
	flag.StringVar(&outputFileName, "out-file", "", "The path to the file to write the final data to")
	flag.BoolVar(&ignoreDuplicates, "ignore-duplicates", false, "Whether or not to ignore duplicate keys")
	flag.Parse()

	if inputDirectory == "" {
		fmt.Println("Source directory not found")
		os.Exit(10)
	}
	if outputFileName == "" {
		fmt.Println("Output file be empty")
		os.Exit(1)
	}

	channelForFilesToParse = make(chan string, 1000)
	channelForDataElementsToSave = make(chan DataElement, 1000)
	abs, err := filepath.Abs(inputDirectory)
	if err != nil {
		log.Fatalf("%s : while reading absolute path for %s", err, inputDirectory)
	}
	err = ReadDirectory(abs)
	if err != nil {
		log.Fatalf("%s : while reading input directory %s", err, inputDirectory)
	}
	wg := sync.WaitGroup{}
	wg.Add(Concurrency)
	if len(channelForFilesToParse) > 0 {
		for i := 0; i < Concurrency; i += 1 {
			go func() {
				if len(channelForFilesToParse) == 0 {
					wg.Done()
					return
				}
				for pathToFileToParse := range channelForFilesToParse {
					err = Parse(pathToFileToParse)
					if err != nil {
						log.Fatalf("%s : while parsing %s", err, pathToFileToParse)
					}
					if len(channelForFilesToParse) == 0 {
						break
					}
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
	Output = make(map[int]DataElement, 0)
	if len(channelForDataElementsToSave) > 0 {
		for de := range channelForDataElementsToSave {
			if 0 == len(channelForDataElementsToSave) {
				break
			}
			_, found := Output[de.ID]
			if found {
				if ignoreDuplicates {
					continue
				}
				fmt.Println("Duplicate data found")
				os.Exit(20)
			}
			Output[de.ID] = de
		}
	}
	outputSlice := make([]DataElement, 0)
	for _, v := range Output {
		outputSlice = append(outputSlice, v)
	}
	sort.Slice(outputSlice, func(i, j int) bool {
		return outputSlice[i].ID < outputSlice[j].ID
	})
	payload, err := json.MarshalIndent(OutputDataStructure{Data: outputSlice}, "", "  ")
	if err != nil {
		log.Fatalf("%s : while marshaling output data to json", err)
	}
	err = ioutil.WriteFile(outputFileName, payload, 0644)
	if err != nil {
		log.Fatalf("%s : while writing output into %s", err, outputFileName)
	}
	os.Exit(0)
}
