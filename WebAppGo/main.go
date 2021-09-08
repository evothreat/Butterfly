package main

import (
	"WebAppGo/api"
)

func main() {
	api.SetupDatabase("root:root@tcp(localhost:3306)/data")
	api.AddTestData()
}
