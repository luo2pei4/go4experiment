package main

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	name := "/home/minio/Documents/chmodtest"
	if err := os.Chmod(name, 0400); err != nil {
		fmt.Printf("Chmod error: %s", err.Error())
	}
	if err := unix.Access(name, unix.R_OK|unix.W_OK); err != nil {
		if errors.Is(err, os.ErrPermission) {
			fmt.Printf("Access error: %s", err.Error())
		}
	} else {
		fmt.Println("Access success.")
	}
}
