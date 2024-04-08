package symmetric

import (
    "bytes"
    "encoding/binary"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "time"

    "github.com/jacastromad/cifra/keys"
)


const NonceLen = 12


// Generate a unique 12-byte nonce using a timestamp and a random value
func generateNonce() []byte {
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


// Encrypts data using GCM. Returns salt+nonce+encrypted_data
func EncGCM(data, pass []byte) ([]byte, error) {
    nonce := generateNonce()  // Number used once

    key := keys.NewKey(pass)

    cBlock, err := aes.NewCipher(key.Bytes)
    if err != nil {
        return nil, err
    }

    aesgcm, err := cipher.NewGCM(cBlock)
    if err != nil {
        return nil, err
    }

    encData := aesgcm.Seal(nil, nonce, data, nil)

    return bytes.Join([][]byte{key.Salt, nonce, encData}, nil), nil
}


// Decrypts the slice of bytes encrypted with EncGCM.
// Expects salt+nonce+encrypted_data
func DecGCM(fields, pass []byte) ([]byte, error) {
    salt := fields[:keys.SaltLen]
    nonce := fields[keys.SaltLen:keys.SaltLen+NonceLen]
    encData := fields[keys.SaltLen+NonceLen:]

    key := keys.GenKey(pass, salt)

    cBlock, err := aes.NewCipher(key.Bytes)
    if err != nil {
        return nil, err
    }

    aesgcm, err := cipher.NewGCM(cBlock)
    if err != nil {
        return nil, err
    }

    data, err := aesgcm.Open(nil, nonce, encData, nil)
    if err != nil {
        return nil, err
    }

    return data, nil
}

// Encrypts data using CFB. Returns salt+iv+encrypted_data
func EncCFB(data, pass []byte) ([]byte, error) {
    iv := make([]byte, aes.BlockSize)
    _, err := rand.Read(iv)
    if err != nil {
        panic("Can't generate a random IV: " + err.Error())
    }

    key := keys.NewKey(pass)

    cBlock, err := aes.NewCipher(key.Bytes)
    if err != nil {
        return nil, err
    }

    aescfb := cipher.NewCFBEncrypter(cBlock, iv)
    encData := make([]byte, len(data))
    aescfb.XORKeyStream(encData, data)

    return bytes.Join([][]byte{key.Salt, iv, encData}, nil), nil
}


// Decrypts the slice of bytes encrypted with EncCFB.
// Expects salt+iv+encrypted_data
func DecCFB(fields, pass []byte) ([]byte, error) {
    salt := fields[:keys.SaltLen]
    iv := fields[keys.SaltLen:keys.SaltLen+aes.BlockSize]
    encData := fields[keys.SaltLen+aes.BlockSize:]

    key := keys.GenKey(pass, salt)

    cBlock, err := aes.NewCipher(key.Bytes)
    if err != nil {
        return nil, err
    }

    aescfb := cipher.NewCFBDecrypter(cBlock, iv)
    data := make([]byte, len(fields)-(keys.SaltLen+aes.BlockSize))
    aescfb.XORKeyStream(data, encData)

    return data, nil
}

