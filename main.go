package main 

import (
	"os"
	"os/signal"
	"fmt"
	"syscall"

	"github.com/MilaSnetkova/TODO-list/internal/app"

)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
	}() 

	a, err := app.New()
	if err != nil {
		fmt.Println("failed to initialize app", err)
		os.Exit(1)
	}

	if err = a.Run(); err != nil {
		defer os.Exit(1)
		fmt.Println("failed to run app", err)
	}
}