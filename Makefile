# expected behaviour
start:
	go run main.go --source-dir=./input/ --ignore-duplicates --out-file=output.json
	cat output.json

# should exit with 20 status code and message `Duplicate data found`
no_dups:
	go run main.go --source-dir=./input/ --out-file=output.json

# should exit with 10 status code and message `Source directory not found`
not_found:
	go run main.go --source-dir=/path/not/found --out-file=output.json

# pack source code into zip
pack:
	git archive --format=zip -o ~/test-data-loader.zip HEAD

