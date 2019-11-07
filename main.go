package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/baunes/api-gatherer/controller"
	"github.com/baunes/api-gatherer/db"
	"github.com/baunes/api-gatherer/db/common"
	"github.com/baunes/api-gatherer/gatherer"
)

type configHTTP struct {
	url string
}

type configDatabase struct {
	host       string
	port       string
	user       string
	pass       string
	database   string
	collection string
}

var httpConfig configHTTP
var databaseConfig configDatabase

func init() {
	flag.StringVar(&httpConfig.url, "url", "", "* The url to gather. Required")
	flag.StringVar(&databaseConfig.host, "db.host", "localhost", "The host of the Mongodb database")
	flag.StringVar(&databaseConfig.port, "db.port", "27017", "The port of the Mongodb database")
	flag.StringVar(&databaseConfig.user, "db.username", " ", "The username of the Mongodb database")
	flag.StringVar(&databaseConfig.pass, "db.password", " ", "The passwrod of the Mongodb database")
	flag.StringVar(&databaseConfig.database, "db.database", "", "* The database to store the requests. Required")
	flag.StringVar(&databaseConfig.collection, "db.collection", "", "* The collection to store the requests. Required")
}

func main() {
	defer exitIfPanic()
	checkArguments()
	err := do(httpConfig, databaseConfig)
	if err != nil {
		log.Printf("Error: %s", err)
	}
}

func checkArguments() {
	flag.Parse()
	databaseConfig.user = strings.TrimSpace(databaseConfig.user)
	databaseConfig.pass = strings.TrimSpace(databaseConfig.pass)
	check(len(httpConfig.url) > 0, "parameter -url is required")
	check(len(databaseConfig.database) > 0, "parameter -db.databse is required")
	check(len(databaseConfig.collection) > 0, "parameter -db.collection is required")
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

func do(cfgHTTP configHTTP, cfgDatabase configDatabase) error {
	log.Printf("Creating Http Client")
	client := gatherer.NewClient()
	log.Printf("Creating Database")
	mongoDB, err := newDatabaseClient(cfgDatabase)
	if err != nil {
		return err
	}
	log.Printf("Creating Repository")
	repository := db.NewGenericRepository(mongoDB.Database(cfgDatabase.database), cfgDatabase.collection)
	log.Printf("Creating Controller")
	ctrl := controller.NewController(client, repository)
	log.Printf("Running controller")
	ctrl.GatherAndSaveURL(cfgHTTP.url)
	return nil
}

func newDatabaseClient(cfgDatabase configDatabase) (common.ClientHelper, error) {
	config := &common.Config{}
	config.Host = cfgDatabase.host
	config.Port = cfgDatabase.port
	if len(cfgDatabase.user) > 0 || len(cfgDatabase.pass) > 0 {
		config.Username = cfgDatabase.user
		config.Password = cfgDatabase.pass
	}
	log.Printf("Creating Database client")
	client, err := common.NewClient(config)
	if err != nil {
		return nil, err
	}
	log.Printf("Connecting to Database")
	err = client.Connect()
	if err != nil {
		return nil, err
	}
	log.Printf("Pinging Database")
	err = client.Ping()
	if err != nil {
		return nil, err
	}
	log.Print("Connected to Database")

	return client, nil
}
