package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"v/internal/server"
)

func main() {

	packagePath := "github.com/pion/webrtc/v3"

	// Execute the `go get` command
	cmd := exec.Command("go", "get", packagePath)

	// Set output to os.Stdout to see output in the terminal
	cmd.Stdout = os.Stdout

	// Set error output to os.Stderr to see errors in the terminal
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running go get: %v\n", err)
		return
	}

	fmt.Println("Package successfully downloaded:", packagePath)





	
	if err := server.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
