package main

import (
	"flag"
	"fmt"
	"github.com/Ezzy77/audioScript-server/internal"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/transcribe"
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
	config           config
	awsConfig        *internal.AwsConfigurations
	supabaseConfig   *internal.SupabaseConfig
	logger           *log.Logger
	s3Client         *s3.Client
	transcribeClient *transcribe.Client
}

func main() {

	awsConfig, err := internal.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v", err)
		os.Exit(1)
	}

	supabaseConfig, err := internal.LoadSupabaseConfig()
	if err != nil {
		fmt.Printf("Error loading Supabase config: %v", err)
		os.Exit(1)
	}

	// creating aws s3 and transcribeClient
	s3Client, transcribeClient, err := internal.InitAwsServices(
		awsConfig.AccessKey,
		awsConfig.SecretKey,
		awsConfig.Region,
	)

	if err != nil {
		fmt.Printf("Error initializing AWS services: %v", err)
		os.Exit(1) // Or however you wish to handle this error
	}

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// init new logger that writes messages to the standard out stream
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config:           cfg,
		logger:           logger,
		awsConfig:        awsConfig,
		supabaseConfig:   supabaseConfig,
		s3Client:         s3Client,
		transcribeClient: transcribeClient,
	}

	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.port),
		Handler:     app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		//WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, server.Addr)
	err = server.ListenAndServe()
	logger.Fatal(err)
}
