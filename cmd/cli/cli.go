package cli

import (
	"fmt"
	"os"

	"github.com/marcuslira/ohploy/lib"
)

func Cli() {
	config, err := lib.LoadConfigFile()
	if err != nil {
		fmt.Printf("ohploy: Error - loading a config file: %v\n", err)
		os.Exit(1)
	}

	mgmt, _ := lib.NewContainerMgmt(config)

	fmt.Printf("ohploy: deploying container from: %s...\n", config.Container.ImageName)

	err = mgmt.DeployContainer()
	if err != nil {
		fmt.Printf("ohploy: Error - Deploying container: %v\n", err)
	}

	fmt.Println("ohploy: Done.")
}
