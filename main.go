package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/baunes/api-gatherer/controller"
	"github.com/baunes/api-gatherer/db"
	"github.com/baunes/api-gatherer/db/common"
	"github.com/baunes/api-gatherer/gatherer"
	"github.com/robfig/cron/v3"
)

type configHTTP struct {
	url  string
	cron string
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
	flag.StringVar(&httpConfig.url, "request.url", "", "* The url to gather. Required")
	flag.StringVar(&httpConfig.cron, "request.cron", "0/10 * * ? * *", "The cron expression. Uses Quartz format http://www.quartz-scheduler.org/documentation/quartz-2.0.2/tutorials/tutorial-lesson-06.html")
	flag.StringVar(&databaseConfig.host, "db.host", "localhost", "The host of the Mongodb database")
	flag.StringVar(&databaseConfig.port, "db.port", "27017", "The port of the Mongodb database")
	flag.StringVar(&databaseConfig.user, "db.username", " ", "The username of the Mongodb database")
	flag.StringVar(&databaseConfig.pass, "db.password", " ", "The passwrod of the Mongodb database")
	flag.StringVar(&databaseConfig.database, "db.database", "", "* The database to store the requests. Required")
	flag.StringVar(&databaseConfig.collection, "db.collection", "", "* The collection to store the requests. Required")
}

type myLogger struct {
}

func (logger myLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Printf(msg, keysAndValues...)
}
func (logger myLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	log.Fatalf(msg, keysAndValues...)
}

func main() {
	defer exitIfPanic()
	checkArguments()
	err := do(newCron(), httpConfig, databaseConfig)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	for {
		time.Sleep(time.Minute)
	}
}

func checkArguments() {
	flag.Parse()
	databaseConfig.user = strings.TrimSpace(databaseConfig.user)
	databaseConfig.pass = strings.TrimSpace(databaseConfig.pass)
	check(len(httpConfig.url) > 0, "parameter -request.url is required")
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

func schedule(c *cron.Cron, spec string, cmd func()) error {
	_, err := c.AddFunc(spec, cmd)
	if err != nil {
		return err
	}
	c.Start()
	return nil
}

func initializeController(cfgHTTP configHTTP, cfgDatabase configDatabase) (controller.Controller, error) {
	log.Printf("Creating Http Client")
	client := gatherer.NewClient()
	log.Printf("Creating Database")
	mongoDB, err := newDatabaseClient(cfgDatabase)
	if err != nil {
		log.Printf("Error: %s", err)
		return nil, err
	}
	log.Printf("Creating Repository")
	repository := db.NewGenericRepository(mongoDB.Database(cfgDatabase.database), cfgDatabase.collection)
	log.Printf("Creating Controller")
	ctrl := controller.NewController(client, repository)
	return ctrl, err
}

func do(c *cron.Cron, cfgHTTP configHTTP, cfgDatabase configDatabase) error {
	ctrl, err := initializeController(cfgHTTP, cfgDatabase)
	if err != nil {
		log.Fatalf("Error initializing controller: %s", err.Error())
	}
	cmd := func() {
		ctrl.GatherAndSaveURL(cfgHTTP.url)
	}
	return schedule(c, cfgHTTP.cron, cmd)
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

func newCron() *cron.Cron {
	return cron.New(cron.WithSeconds(), cron.WithLogger(myLogger{}))
}
