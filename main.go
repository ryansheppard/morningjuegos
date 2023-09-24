package main

import (
	"fmt"
	"os"

	"github.com/ryansheppard/morningjuegos/cmd"
	"github.com/ryansheppard/morningjuegos/internal/database"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	err := database.CreateConnection(dbPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd.Execute()
}
