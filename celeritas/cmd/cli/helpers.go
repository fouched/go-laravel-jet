package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"os"
)

func setup() {
	err := godotenv.Load()
	if err != nil {
		exitGracefully(err)
	}

	path, err := os.Getwd()
	if err != nil {
		exitGracefully(err)
	}

	cel.RootPath = path
	cel.DB.Type = os.Getenv("DATABASE_TYPE")
}

func getDSN() string {
	dbType := cel.DB.Type

	//we use pgx, but golang migrate uses different driver
	//so convert dsn to work with golang migrate
	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		}
		return dsn
	}
	return "mysql://" + cel.BuildDSN()
}

func showHelp() {
	color.Yellow(`Available commands:

    help                     - show help
    version                  - print version
    make auth                - creates authentication tables, models and middleware
    make handler <name>      - creates a stub handler in the handlers directory
    make key                 - creates a random 32 character encryption key 
    make model <name>        - creates a new model in the data directory
    make session             - creates a new table as a session store
    
    make migration <name>    - creates new up and down migrations
    migrate                  - runs all up migrations
    migrate down             - reverses most recent migration
    migrate reset            - runs all down migrations, then all up migrations
    `)
}
