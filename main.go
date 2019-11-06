package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
)

var urlToGather string
var hostMongo string
var portMongo string
var userMongo string
var passMongo string
var databseMongo string
var collectionMongo string

func init() {
	flag.StringVar(&urlToGather, "url", "", "* The url to gather. Required")
	flag.StringVar(&hostMongo, "db.host", "localhost", "The host of the Mongodb database")
	flag.StringVar(&portMongo, "db.port", "27017", "The port of the Mongodb database")
	flag.StringVar(&userMongo, "db.username", " ", "The username of the Mongodb database")
	flag.StringVar(&passMongo, "db.password", " ", "The passwrod of the Mongodb database")
	flag.StringVar(&databseMongo, "db.database", "", "* The database to store the requests. Required")
	flag.StringVar(&collectionMongo, "db.collection", "", "* The collection to store the requests. Required")
}

func main() {
	defer exitIfPanic()
	checkArguments()
	do()
}

func checkArguments() {
	flag.Parse()
	check(len(urlToGather) > 0, "parameter -url is required")
	check(len(databseMongo) > 0, "parameter -db.databse is required")
	check(len(collectionMongo) > 0, "parameter -db.collection is required")
}

func exitIfPanic() {
	if message := recover(); message != nil {
		fmt.Fprintf(os.Stderr, "%s\n", message)
		syscall.Exit(2)
	}
}

func check(result bool, message string) {
	if !result {
		panic(fmt.Sprintf("error: %s", message))
	}
}

func do() {
}
