package commands

import (
	"fmt"
	"os"
)

func Stop() {
	fmt.Print("The stop is complete")
	os.Exit(0)
}
