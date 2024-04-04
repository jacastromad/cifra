package keys


import (
    "crypto/rand"
    "crypto/sha512"
    "golang.org/x/crypto/pbkdf2"
)


const (
    Iter = 500000  // Iteration count
    KeyLen = 32    // AES 256 - 32 bytes
    SaltLen = 64   // Salt length
)


// Contains a derived key
type key struct {
    Bytes []byte
    Salt []byte
}


// Key constructor. Returns a key struct with a random salt
func NewKey(pass []byte) *key {
    // Generate saltLen random bytes
    saltb := make([]byte, SaltLen)
    _, err := rand.Read(saltb)
    if err != nil {
        panic("Can't generate a random salt for a new key: " + err.Error())
    }

    keyb := pbkdf2.Key(pass, saltb, Iter, KeyLen, sha512.New)

    k := key {
        Bytes: keyb,
        Salt: saltb,
    }

    return &k
}


// Key constructor. Returns key struct with the given salt
func GenKey(pass, salt []byte) *key {
    keyb := pbkdf2.Key(pass, salt, Iter, KeyLen, sha512.New)

    k := key {
        Bytes: keyb,
        Salt: salt,
    }

    return &k
}

