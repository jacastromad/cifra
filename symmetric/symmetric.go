package symmetric

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/binary"
    "errors"
    "fmt"
    "time"

    "github.com/jacastromad/cifra/keys"
)


const NonceLen = 12


// Generate a unique 12-byte nonce using a timestamp and a random value
func generateNonce() ([]byte, error) {
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
        return nil, err
    }

    // Combine the timestamp and random bytes to form a 12-byte nonce
    nonce := append(timestampBytes, randomBytes...)

    return nonce, nil
}


// Encrypts data using GCM. Returns salt+nonce+encrypted_data
func EncGCM(data, pass []byte) ([]byte, error) {
    nonce, err := generateNonce()  // Number used once
    if err != nil {
        return nil, fmt.Errorf("EncGCM: error generating nonce: %w", err)
    }

    key, err := keys.NewKey(pass)
    if err != nil {
        return nil, fmt.Errorf("EncGCM: error generating key: %w", err)
    }

    cBlock, err := aes.NewCipher(key.Bytes)
    if err != nil {
        return nil, fmt.Errorf("EncGCM: error creating cipher.Block: %w", err)
    }

    aesgcm, err := cipher.NewGCM(cBlock)
    if err != nil {
        return nil, fmt.Errorf("EncGCM: error wrapping cipher in GCM: %w", err)
    }

    encData := aesgcm.Seal(nil, nonce, data, nil)

    return bytes.Join([][]byte{key.Salt, nonce, encData}, nil), nil
}


// Decrypts the slice of bytes encrypted with EncGCM.
// Expects salt+nonce+encrypted_data
func DecGCM(fields, pass []byte) ([]byte, error) {
    if len(fields) < keys.SaltLen+NonceLen {
        return nil, errors.New("Corrupted data")
    }
    salt := fields[:keys.SaltLen]
    nonce := fields[keys.SaltLen:keys.SaltLen+NonceLen]
    encData := fields[keys.SaltLen+NonceLen:]

    key := keys.GenKey(pass, salt)

    cBlock, err := aes.NewCipher(key.Bytes)
    if err != nil {
        return nil, fmt.Errorf("DecGCM: error creating cipher.Block: %w", err)
    }

    aesgcm, err := cipher.NewGCM(cBlock)
    if err != nil {
        return nil, fmt.Errorf("DecGCM: error wrapping cipher in GCM: %w", err)
    }

    data, err := aesgcm.Open(nil, nonce, encData, nil)
    if err != nil {
        return nil, fmt.Errorf("DecGCM: error decypting ciphertext: %w", err)
    }

    return data, nil
}

// Encrypts data using CFB. Returns salt+iv+encrypted_data
func EncCFB(data, pass []byte) ([]byte, error) {
    iv := make([]byte, aes.BlockSize)
    _, err := rand.Read(iv)
    if err != nil {
        return nil, fmt.Errorf("EncCFB: error creating random IV: %w", err)
    }

    key, err := keys.NewKey(pass)
    if err != nil {
        return nil, fmt.Errorf("EncCFB: error generating key: %w", err)
    }

    cBlock, err := aes.NewCipher(key.Bytes)
    if err != nil {
        return nil, fmt.Errorf("EncCFB: error creating cipher.Block: %w", err)
    }

    aescfb := cipher.NewCFBEncrypter(cBlock, iv)
    encData := make([]byte, len(data))
    aescfb.XORKeyStream(encData, data)

    return bytes.Join([][]byte{key.Salt, iv, encData}, nil), nil
}


// Decrypts the slice of bytes encrypted with EncCFB.
// Expects salt+iv+encrypted_data
func DecCFB(fields, pass []byte) ([]byte, error) {
    if len(fields) < keys.SaltLen+aes.BlockSize {
        return nil, errors.New("Corrupted data")
    }
    salt := fields[:keys.SaltLen]
    iv := fields[keys.SaltLen:keys.SaltLen+aes.BlockSize]
    encData := fields[keys.SaltLen+aes.BlockSize:]

    key := keys.GenKey(pass, salt)

    cBlock, err := aes.NewCipher(key.Bytes)
    if err != nil {
        return nil, fmt.Errorf("DecCFB: error creating cipher.Block: %w", err)
    }

    aescfb := cipher.NewCFBDecrypter(cBlock, iv)
    data := make([]byte, len(fields)-(keys.SaltLen+aes.BlockSize))
    aescfb.XORKeyStream(data, encData)

    return data, nil
}

