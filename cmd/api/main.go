package main

import (
	"flag"
	"fmt"
	"github.com/Ezzy77/transcribe-go/internal"
	"log"
	"net/http"
	"os"
	"time"
)

//const version = "1.0.0"

// holds application config
type config struct {
	port int
	env  string
}

// application struct holds dependencies for http
// handler and middleware
type application struct {
	config    config
	awsConfig *internal.AwsConfigurations
	logger    *log.Logger
}

func main() {

	awsConfig, err := internal.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v", err)
		os.Exit(1) // Or however you wish to handle this error
	}
	fmt.Println("here", len(awsConfig.AWS.AccessKey))
	fmt.Println(awsConfig.AWS.SecretKey)
	fmt.Println(awsConfig.AWS.Region)

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// init new logger that writes messages to the standard out stream
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config:    cfg,
		logger:    logger,
		awsConfig: awsConfig,
		////store:    store,
		//users:    &models.UserModel{DB: db},
		//expenses: &models.ExpenseModel{DB: db},
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, server.Addr)
	err = server.ListenAndServe()
	logger.Fatal(err)
}
