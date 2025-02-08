package main

import "os"

func main() {
	defer os.Exit(0) // want "os.Exit used in main function of main package"

	other()

	if true {
		os.Exit(1) // want "os.Exit used in main function of main package"
	}

	os.Exit(1) // want "os.Exit used in main function of main package"
}

func other() {
	os.Exit(1)
}
