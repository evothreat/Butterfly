package main

import "fmt"

func main() {
	var w Worker = Worker{
		id: "azazaza",
	}
	fmt.Printf("Hello World %s\n", w.id)
}
