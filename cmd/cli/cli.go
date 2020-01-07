package cli

import (
	"fmt"

	"github.com/marcuslira/ohploy/lib"
)

func Cli() {
	config, _ := lib.NewEnvConfig()
	mgmt, _ := lib.NewContainerMgmt(config)

	fmt.Println("ohploy: deploying container...")
	err := mgmt.DeployContainer()
	if err != nil {
		fmt.Printf("ohploy: Error - Deploying container: %v\n", err)
	}

	fmt.Println("ohploy: Done.")
}
