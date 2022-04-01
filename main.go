package main

import (
	"fmt"
	"github.com/gcamlicali/RESTHW/app"
	"log"
)

func main() {

	exec := &app.App{}
	err := exec.Initializer()

	if err != nil {
		log.Fatal("Application not initialized")
	}

	fmt.Println("Server Ready to UP")
	exec.Run(":3000")
}
