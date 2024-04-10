package utils


import (
    "os"
    "io/fs"
    "encoding/base64"
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
        panic("Error reading file: " + err.Error())
    }

    if b64 {
        b64str := string(rdata)
        return base64.StdEncoding.DecodeString(b64str)
    } else {
        return rdata, nil
    }
}

