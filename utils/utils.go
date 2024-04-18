package utils

import (
    "encoding/base64"
    "fmt"
    "io/fs"
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


// Write byte slice as is or encoded as base64
func WriteFile(filename string, data []byte, perm fs.FileMode, b64 bool) error {
    var err error

    if b64 {
        b64str := base64.StdEncoding.EncodeToString(data)
        err = os.WriteFile(filename, []byte(b64str), perm)
    } else {
        err = os.WriteFile(filename, data, perm)
    }

    return err
}


// Read file as byte slice or encoded as base64
func ReadFile(filename string, b64 bool) ([]byte, error) {
    rdata, err := os.ReadFile(filename)

    if err != nil {
        return nil, fmt.Errorf("ReadFile: error reading file: %w", err)
    }

    if b64 {
        b64str := string(rdata)
        decbytes, err := base64.StdEncoding.DecodeString(b64str)
        if err != nil {
            return nil, fmt.Errorf("ReadFile: error decoding base64: %w", err)
        } else {
            return decbytes, nil
        }
    } else {
        return rdata, nil
    }
}

