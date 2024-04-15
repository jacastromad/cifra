package keys


import (
    "testing"
    "bytes"
)


func TestNewKey(t *testing.T) {
    passwords := [][]byte{
        []byte("1234"),
        []byte("-.^'ä¿++123abc"),
        []byte(""),
        []byte("Contraseña"),
        []byte(" "),
    }

    for _, pass := range passwords {
        key := NewKey(pass)
        // Test that salt is generated and has the correct length
        if len(key.Salt) != SaltLen {
            t.Errorf("Expected salt length of %d, got %d", SaltLen, len(key.Salt))
        }

        // Test that the key is generated and has the correct length
        if len(key.Bytes) != KeyLen {
            t.Errorf("Expected key length of %d, got %d", KeyLen, len(key.Bytes))
        }
    }
}

func TestGenKey(t *testing.T) {
    plist1 := [][]byte{
        []byte("1234"),
        []byte("-.^'ä¿++123abc"),
        []byte(""),
        []byte("Contraseña"),
        []byte(" "),
    }

    plist2 := [][]byte{
        []byte("1234 "),
        []byte("-.^'á¿++123abc"),
        []byte(" "),
        []byte("contraseña"),
        []byte("  "),
    }

    for i, _ := range plist1 {
        k1 := NewKey(plist1[i])
        k2 := GenKey(plist1[i], k1.Salt)
        k3 := GenKey(plist2[i], k1.Salt)
        // Check that the provided salt is used
        if !bytes.Equal(k1.Salt, k2.Salt) {
            t.Errorf("GenKey did not use provided salt %v", k1.Salt)
        }
        // Same password showld produce the same output
        if !bytes.Equal(k1.Bytes, k2.Bytes) {
            t.Errorf("NewKey != GenKey for password %s.", plist1[i])
        }
        // Different password should produce a different output
        if bytes.Equal(k1.Bytes, k3.Bytes) {
            t.Errorf("NewKey == GenKey for '%s' and '%s'.", plist1[i], plist2[i])
        }
    }
}

