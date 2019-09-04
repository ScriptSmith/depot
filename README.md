# ![](assets/logo.png) Depot

[![Build Status](https://travis-ci.org/ScriptSmith/depot.svg?branch=master)](https://travis-ci.org/ScriptSmith/depot)
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/scriptsmith/depot)](https://hub.docker.com/r/scriptsmith/depot)

A fast and simple HTTP server to store files from scraping and processing jobs.

Depot is designed to be used with an attached disk as a 'quick and dirty' object storage microservice.

[![Run on Google Cloud](https://storage.googleapis.com/cloudrun/button.svg)](https://console.cloud.google.com/cloudshell/editor?shellonly=true&cloudshell_image=gcr.io/cloudrun/button&cloudshell_git_repo=https://github.com/scriptsmith/depot.git)

## Features
- Download files with a GET request
- Upload files with a PUT request
- Download a zipped copy of a job's files


## Build
```
go get -v github.com/scriptsmith/depot
go build -v github.com/scriptsmith/depot
```

## Run
Create a root directory, or pick an existing one:

```
# mkdir /tmp/dr
export DEPOT_ROOT=/tmp/dr
```

Run

```
go run github.com/scriptsmith/depot
```

```
$ ./depot 
2019/07/25 17:13:03 running at http://127.0.0.1:8080
```

Access the page to see instructions and a list of current jobs

## Docker

Run `scriptsmith/depot` and forward port `8080`
```
docker run -p 8080:8080 scriptsmith/depot
```

Use authentication and attach dir on host:

```
docker run -p 8080:8080 \
 -e DEPOT_USER=depot \
 -e DEPOT_PASS=pass \
 -v /tmp/dr:/data \
 scriptsmith/depot
```

