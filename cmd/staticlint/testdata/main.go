package main

import "os"

func main() {
	os.Exit(1) // want "os.Exit used in main function of main package"
}
