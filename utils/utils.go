package utils


import (
    "encoding/binary"
    "crypto/rand"
    "time"
    "os"
)


const NonceLen = 12


// Generate a unique 12-byte nonce using a timestamp and a random value
func GenerateNonce() []byte {
    // Get the current time in nanoseconds since the Unix epoch
    epochNano := time.Now().UnixNano()

    // Convert the timestamp to a byte slice (8 bytes for i64)
    buf := make([]byte, 8)
    binary.BigEndian.PutUint64(buf, uint64(epochNano))

    // Skip the two most significant bytes
    timestampBytes := buf[2:]

    // Generate 6 bytes of random data
    randomBytes := make([]byte, 6)
    if _, err := rand.Read(randomBytes); err != nil {
        panic(err)
    }

    // Combine the timestamp and random bytes to form a 12-byte nonce
    nonce := append(timestampBytes, randomBytes...)

    return nonce
}


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


