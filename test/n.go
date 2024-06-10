package main

import (
	"fmt"

	"github.com/sqweek/dialog"
)

func main() {
	file, err := dialog.File().Title("").Filter("png", "png").Load()
	fmt.Println(file)
	fmt.Println("Error:", err)
	// dialog.Message("You chose file: %s", file).Title("Goodbye world!").Error()
	// dialog.Directory().Title("Now find a dir").Browse()
}