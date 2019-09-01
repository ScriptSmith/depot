# Depot
[![Build Status](https://travis-ci.org/ScriptSmith/depot.svg?branch=master)](https://travis-ci.org/ScriptSmith/depot)

A fast and simple HTTP server to store files from scraping and processing jobs.

Depot is designed to be used with an attached disk as a 'quick and dirty' object storage microservice.

[![Run on Google Cloud](https://storage.googleapis.com/cloudrun/button.svg)](https://console.cloud.google.com/cloudshell/editor?shellonly=true&cloudshell_image=gcr.io/cloudrun/button&cloudshell_git_repo=https://github.com/scriptsmith/depot.git)

## Features
- Download files with a GET request
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

