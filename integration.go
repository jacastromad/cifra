//go:build itest

package main

import (
    "fmt"
)

// Read password used for integration tests
func readPassword(fd int) ([]byte, error) {
    fmt.Print("[WARNING! INSECURE INPUT METHOD]: ")
    var pass []byte
    _, err := fmt.Scanln(&pass)
    return pass, err
}

