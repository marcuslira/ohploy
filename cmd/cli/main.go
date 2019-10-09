package main

import (
	"fmt"

	"github.com/marcuslira/ohploy/lib"
)

func main() {
	fmt.Println("ohploy: deploying container...")
	imageName := "docker.io/marcuslira/aspiratracker:latest"

	err := lib.DeployContainer(imageName)
	if err != nil {
		fmt.Printf("ohploy: Error - Deploying container: %v\n", err)
	}

	fmt.Println("ohploy: Done.")
}
