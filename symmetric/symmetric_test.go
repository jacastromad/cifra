package symmetric

import (
    "bytes"
    "testing"
    "crypto/rand"
)

func TestGenerateNonce(t *testing.T) {
    nonce, err := generateNonce()
    if err != nil {
        t.Errorf("Error generating nonce: %s", err)
    }
    if len(nonce) != NonceLen {
        t.Errorf("Expected nonce to be %d bytes long, got %d", NonceLen, len(nonce))
    }
}

func TestEncGCMAndDecGCM(t *testing.T) {
    passwords := [][]byte{
        []byte("1234"),
        []byte("-.^'ä¿++123abc"),
        []byte(""),
        []byte("Contraseña"),
        []byte(" "),
    }
    lengths := []int{0, 100, 10000}

    for _, pass := range passwords {
        for _, length := range lengths {
            data := make([]byte, length)
            _, err := rand.Read(data)
            if err != nil {
                t.Fatalf("Failed to generate random data: %s", err)
            }
            encData, err := EncGCM(data, pass)
            if err != nil {
                t.Errorf("EncGCM failed: %s", err)
            }
            if bytes.Equal(data, encData) {
                t.Errorf("EncGCM: encrypted data is the same as the plan data")
            }
            decData, err := DecGCM(encData, pass)
            if err != nil {
                t.Errorf("DecGCM failed on pass '%s', len '%d': %s", pass, length, err)
            }
            if !bytes.Equal(data, decData) {
                t.Errorf("DecGCM: decrypted data does not match the original")
            }
        }
    }
}

func TestEncCFBAndDecCFB(t *testing.T) {
    passwords := [][]byte{
        []byte("1234"),
        []byte("-.^'ä¿++123abc"),
        []byte(""),
        []byte("Contraseña"),
        []byte(" "),
    }
    lengths := []int{0, 100, 10000}

    for _, pass := range passwords {
        for _, length := range lengths {
            data := make([]byte, length)
            _, err := rand.Read(data)
            if err != nil {
                t.Fatalf("Failed to generate random data: %s", err)
            }
            encData, err := EncCFB(data, pass)
            if err != nil {
                t.Errorf("EncGCM failed: %s", err)
            }
            if bytes.Equal(data, encData) {
                t.Errorf("EncGCM: encrypted data is the same as the plan data")
            }
            decData, err := DecCFB(encData, pass)
            if err != nil {
                t.Errorf("DecGCM failed on pass '%s', len '%d': %s", pass, length, err)
            }
            if !bytes.Equal(data, decData) {
                t.Errorf("DecGCM: decrypted data does not match the original")
            }
        }
    }
}

func TestFailureModes(t *testing.T) {
    passwords := [][]byte{
        []byte("1234"),
        []byte("-.^'ä¿++1a"),
        []byte(""),
    }

    encDatas := [][]byte{  // Length sould be at least Salt+Nonce || Salt+IV
        {},
        {0, 0, 0, 0, 0},
        {100, 101, 102, 103, 104, 105, 106, 107, 108, 109},
    }

    for _, pass := range passwords {
        for _, encData := range encDatas {
            // Testing DecGCM with corrupted data
            _, err := DecGCM(encData, pass)
            if err == nil {
                t.Error("DecGCM should fail (corrupted data) but no error was returned")
            }

            // Testing DecCFB with corrupted data
            _, err = DecCFB(encData, pass)
            if err == nil {
                t.Error("DecCFB should fail (corrupted data) but no error was returned")
            }

            // Testing DecGCM with wrong password
            data := []byte("Secret data")
            encryptedData, _ := EncGCM(data, pass)
            _, err = DecGCM(encryptedData, []byte("badpassword"))
            if err == nil {
                t.Error("DecGCM should fail (wrong password) but no error was returned")
            }
        }
    }

}

