# Depot
A simple HTTP server to store scraping jobs

## Features
- Access files with a GET request
- Upload files with a PUT request
- Download a zipped copy of a job's files

## Build
```
go build .
```

## Run
Create a root directory, or pick an existing one:

```
# mkdir /tmp/dr
export DEPOT_ROOT=/tmp/dr
```

Run

```
./depot
```

```
$ ./depot 
2019/07/25 17:13:03 running at http://127.0.0.1:8080
```

Access the page to see instructions and a list of current jobs