test-data-loader
========================

See `Go - Data Loader.pdf` for details.


Running
=========================

Ordinary execution.

```shell
$ make start 
```

It should output data without duplicates:
```
go run main.go --source-dir=./input/ --ignore-duplicates --out-file=output.json
cat output.json
{
  "data": [
    {
      "id": 100,
      "first_name": "John",
      "last_name": "Doe"
    },
    {
      "id": 101,
      "first_name": "Jane",
      "last_name": "Doe"
    },
    {
      "id": 102,
      "first_name": "Jackson",
      "last_name": "Doe"
    },
    {
      "id": 103,
      "first_name": "Jack",
      "last_name": "Doe"
    },
    {
      "id": 104,
      "first_name": "Vincent",
      "last_name": "Doe"
    }
  ]
}
```

Duplicate check failed:
```shell
$ make no_dups
```
It should output
```
go run main.go --source-dir=./input/ --out-file=output.json
Duplicate data found
exit status 20
make: *** [Makefile:8: no_dups] Error 1
```

No input directory found

```shell
$ make not_found
```

```
go run main.go --source-dir=/path/not/found --out-file=output.json
Source directory not found
exit status 10
make: *** [Makefile:12: not_found] Error 1
```
