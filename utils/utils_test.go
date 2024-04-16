package utils


import (
    "os"
    "bytes"
    "crypto/rand"
    "testing"
)


const filename = "iui2187679oquowdufh123312397887665.tmp"


func TestFileExists(t *testing.T) {
    // Ensure the file does not exist at the start
    os.Remove(filename)

    if FileExists(filename) {
        t.Errorf("FileExists returns true for non-existing file")
    }

    // Create the file
    _, err := os.Create(filename)
    if err != nil {
        t.Fatal("Error creating test file:", err)
    }
    defer os.Remove(filename)

    if !FileExists(filename) {
        t.Errorf("FileExists returns false for existing file")
    }
}


func TestWriteFileAndReadFile(t *testing.T) {
    defer os.Remove(filename)

    data := make([]byte, 1000)
    _, err := rand.Read(data)
    if err != nil {
        t.Fatalf("Failed to generate random data: %s", err)
    }

    // Test writing random data
    err = WriteFile(filename, data, 0644, false)
    if err != nil {
        t.Fatalf("WriteFile failed: %s", err)
    }

    // Test reading random data
    readData, err := ReadFile(filename, false)
    if err != nil {
        t.Fatalf("ReadFile failed: %s", err)
    }
    if !bytes.Equal(data, readData) {
        t.Errorf("Read data did not match original data")
    }

    // Test writing base64 data
    err = WriteFile(filename, data, 0644, true)
    if err != nil {
        t.Fatalf("WriteFile with base64 failed: %s", err)
    }

    // Test reading base64 data
    readData, err = ReadFile(filename, true)
    if err != nil {
        t.Fatalf("ReadFile with base64 failed: %s", err)
    }
    if !bytes.Equal(data, readData) {
        t.Errorf("Expected %s, got %s", string(data), string(readData))
    }
}

//func TestReadFilePanic(t *testing.T) {
//    filename := "nonexistent.txt"
//
//    defer func() {
//        if r := recover(); r == nil {
//            t.Errorf("ReadFile did not panic on missing file")
//        }
//    }()
//
//    _, _ = ReadFile(filename, false)
//}

