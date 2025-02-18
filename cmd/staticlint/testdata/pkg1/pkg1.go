package pkg1

import "os"

// Another - просто функция в тестовов пакете.
func Another() {
	defer os.Exit(0)
}

func other() {
	os.Exit(1)
}
