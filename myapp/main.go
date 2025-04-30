package main

import (
	"context"
	"fmt"
	"github.com/fouched/celeritas"
	"myapp/data"
	"myapp/handlers"
	"myapp/middleware"
	"os"
	"os/signal"
	"syscall"
)

type application struct {
	App        *celeritas.Celeritas
	Handlers   *handlers.Handlers
	Models     data.Models
	Middleware *middleware.Middleware
}

func main() {
	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the application
	go func() {
		c := initApplication()
		c.App.ListenAndServe()
	}()

	// Listen for the interrupt signal to complete, giving defer statements time to run
	<-ctx.Done()
	fmt.Println("Graceful shutdown complete")
}
