package main 

import (
	"os"
	"fmt"

	"github.com/MilaSnetkova/TODO-list/internal/app"

)

func main() {
	a, err := app.New()
	if err != nil {
		fmt.Println("failed to initialize app", err)
		os.Exit(1)
	}

	if err = a.Run(); err != nil {
		fmt.Println("failed to run app", err)
	}
}