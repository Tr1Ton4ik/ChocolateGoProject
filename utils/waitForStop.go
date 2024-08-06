package utils

import (
	"fmt"
	"os"
	"strings"
)

func WaitForStop() {
	var stop string
	for {
		fmt.Scan(&stop)
		if strings.ToLower(stop) == "s" {
			fmt.Println("Successfully stopped")
			os.Exit(0)
		}
	}
}
