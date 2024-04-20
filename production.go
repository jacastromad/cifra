//go:build !itest

package main

import "golang.org/x/term"

// Read password without local echo
func readPassword(fd int) ([]byte, error) {
    return term.ReadPassword(fd)
}

