# Simple Invert Index implement by golang
Simple invert index implemention, Use the Gin framework to do the RESTful API. In here provide Add and Search operation. 

## How it work


### add to invert index
The Add API will input the filename which need to append to the invert index
this file should have the particular restriction, when new the invert index server, User can input the parameter as the search key word. for example
if use the "ID", "Name", "title" as the key word.
The file header should like
```
ID:`value`
Name:`value`
title:`value`
```
then will use the key word and timestamp to hash as the document ID and all the key word is the invert index key the document ID is value.
finally rename the input file name to document for easy to read.

### search from invert index
input the key word then will return the document ID or nothing.


## Install
```shell
go get github.com/Noahnut/invertIndex
```

## Usage
1. create the invertIndex struct
2. New the server and input the what search word want to read from the file

```go
in := invert.InvertIndex{}
in.NewInvertIndex("ID", "Name", "title")
```

### **/AddNewDocument**  
**POST**

### **/GetDocument**
**GET**
