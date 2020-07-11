package main

import (
	"flag"
	"fmt"
)

func main() {
	name := flag.String("name", "Paweł", "")
	flag.Parse()
	fmt.Printf("Cześć, %s\n", *name)
}
