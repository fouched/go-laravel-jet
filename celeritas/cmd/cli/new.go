package main

import (
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
	"strings"
)

func doNew(appName string) {
	appName = strings.ToLower(appName)

	// sanitize appName
	if strings.Contains(appName, "/") {
		exploded := strings.SplitAfter(appName, "/")
		appName = exploded[len(exploded)-1]
	}
	log.Println("App name is", appName)

	// git clone skeleton application
	color.Green("\tCloning repository")
	_, err := git.PlainClone("./"+appName, false, &git.CloneOptions{
		URL:      "https://github.com/fouched/celeritas-app.git",
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		exitGracefully(err)
	}

	// remove .git directory

	// create a ready to go .env file

	// create a makefile

	// update go.mod file

	// update existing .go files with correct name/imports

	// run go mod tidy in project directory
}
