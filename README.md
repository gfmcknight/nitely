# Nitely

Nitely is a local build service that allows users to build their
projects at specified times more easily.

## Prerequisites

Nitely currently requires Golang and an SQLite3 driver from mattn.
Run
```
go get github.com/mattn/go-sqlite3
```
## Installing

Once you have cloned the repo you can install it with:
```
go install
```
and then
```
go build
```
You will need to set the environment variable NitelyPath to the install
location of Nitely, as well as add this location to the PATH.

## Usage
More arguments and documentation will come
```
Nitely add [-n <Build Name>] [-d <Directory>] [-b <Branch Name>]
Nitely add -s <Service Filename> [-d <Directory>] [-a <Args>] [-n <Service Name>]
Nitely remove <Name> [-s]
Nitely build <Build Name>
Nitely start <Service Name>
Nitely list <builds/services>
```

## Roadmap

* Support for starting services/processes that a project requires.
* Support for fetching and inflating a remote reference.
* Reading test results from a file and storing them in the database.
* Set the build passing/failing status for projects on GitHub.