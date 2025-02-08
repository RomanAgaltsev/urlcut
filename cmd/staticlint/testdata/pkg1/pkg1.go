package pkg1

import "os"

func Another() {
	defer os.Exit(0)
}

func other() {
	os.Exit(1)
}
