package utils


import (
    "os"
)


// Returns true if file exists
func FileExists(filename string) bool {
    _, err := os.Stat(filename)
    if err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}


