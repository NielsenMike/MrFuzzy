package main

import (
	"fmt"
	"github.com/mholt/archiver/v3"
	"os"
)

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	err = archiver.Unarchive(workingDir + "/wep-app-node.tar", workingDir + "/wep-app-node")
	fmt.Println(err)
}

