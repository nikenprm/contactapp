package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/contactapp/repository"
)

func executeMigration(args []string) {
	if len(args) < 2 {
		fmt.Printf("Missing migration command\n")
		os.Exit(1)
	}
	if args[1] == "up" || args[1] == "down" {
		var (
			num int
			err error
		)
		if len(args) > 2 {
			num, err = strconv.Atoi(args[2])
			if err != nil {
				fmt.Printf("Invalid migration parameter: %s\n", err.Error())
				os.Exit(1)
			}
		}
		if args[1] == "down" {
			if num == 0 {
				num = -1
			} else {
				num = num * -1
			}
		}
		if errors := repository.Migrate(*migrationPath, num); len(errors) > 0 {
			for i := range errors {
				fmt.Printf("Migration error: %s\n", errors[i].Error())
				os.Exit(1)
			}
		}
	} else if args[1] == "generate" {
		if len(args) < 3 {
			fmt.Printf("Missing migration file name\n")
			os.Exit(1)
		}
		version := strconv.FormatInt(time.Now().UTC().Unix(), 10)
		directions := [2]string{"up", "down"}
		for i := range directions {
			file, err := os.Create(fmt.Sprintf("%s/%s_%s.%s.sql", *migrationPath, version, args[2], directions[i]))
			if err != nil {
				fmt.Printf("Error creating migration file: %s\n", err.Error())
				os.Exit(1)
			}
			file.Close()
		}
	} else {
		fmt.Printf("Unknown migration command: %s\n", args[1])
	}
}
