package main

import "SFC-Scheduler/pkg/database"

func main() {
	database.InitDB()
	database.Query()
}
