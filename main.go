package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "syscall"

    "github.com/jacastromad/cifra/symmetric"
    "github.com/jacastromad/cifra/utils"
    "golang.org/x/term"
)


var binName = filepath.Base(os.Args[0])
var helpStr = fmt.Sprintf(`Symmetric encryption

Usage:
    %s [flags] filename

When the output file name is not specified using the -o flag, the program
automatically appends the suffix '.cif' for encryption operations or '.dec'
for decryption operations to the input file name, provided the name is not
already taken.

Flags:
`, binName)


// Print commands help
func printHelp() {
    flag.CommandLine.SetOutput(os.Stderr)
    fmt.Fprintf(os.Stderr, helpStr)
    flag.PrintDefaults()
    fmt.Println()
}


func main() {
    // Available modes
    var gcm bool
    var cfb bool
    // Decrypt
    var dec bool
    // Input and output filenames
    var iFilename string
    var oFilename string
    // Read/Write base64
    var b64 bool
    // Print help
    var help bool

    //flag.CommandLine.SetOutput(io.Discard)
    flag.BoolVar(&gcm, "gcm", false, "Galois/Counter Mode")
    flag.BoolVar(&cfb, "cfb", false, "cipher feedback mode")
    flag.BoolVar(&dec, "dec", false, "Decrypt input file")
    flag.BoolVar(&b64, "b64", false, "Read/Write data encoded as base64")
    flag.StringVar(&oFilename, "o", "", "Output file name")
    flag.BoolVar(&help, "help", false, "Print help")
    flag.Usage = printHelp

    flag.Parse()

    if help {
        printHelp()
        os.Exit(1)
    }
   
    // Begin error handling and default values

    // No input filename
    args := flag.Args()
    if len(args) < 1 {
        fmt.Printf("Error: No input filename.\n")
        fmt.Printf("Please use the -help flag to display usage information.\n")
        os.Exit(1)
    }
    iFilename = args[0] // The input filename is the last argument

    // More than one mode selected
    if gcm && cfb {
        fmt.Printf("Error: Multiple modes selected.\n")
        fmt.Printf("Please specify just one mode.\n")
        os.Exit(1)
    }

    // If no mode was selected, gcm is the default
    if !(gcm || cfb) {
        gcm = true
    }

    // Input file does not exist
    if !utils.FileExists(iFilename) {
        fmt.Printf("Error: file \"%s\" does not exist!\n", iFilename)
        os.Exit(1)
    }

    // If no output filename provided
    if oFilename == "" {
        if dec {
            oFilename = iFilename + ".dec"
        } else {
            oFilename = iFilename + ".cif"
        }
    }

    // If output file already exists
    if utils.FileExists(oFilename) {
        fmt.Printf("Error: file \"%s\" already exist!\n", oFilename)
        os.Exit(1)
    }

    // End error handling and default values


    // Ask for password
    fmt.Print("Password: ")
    pass, err := term.ReadPassword(int(syscall.Stdin))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println()

    var data []byte
    if dec {  // Decrypt
        // Read file to decrypt
        bytes, err := utils.ReadFile(iFilename, b64)
        if err != nil {
            panic("Error reading input file: " + err.Error())
        }
        if gcm {  // GCM
            data, err = symmetric.DecGCM(bytes, pass)
        } else {  // CFB
            data, err = symmetric.DecCFB(bytes, pass)
        }
        if err != nil {
            panic("Error while decrypting data: " + err.Error())
        }
        // Write the unencrypted data
        err = utils.WriteFile(oFilename, data, 0600, false) // rw for user only
    } else {  // Encrypt
        // Read file to encrypt
        bytes, err := utils.ReadFile(iFilename, false)
        if err != nil {
            panic("Error reading input file: " + err.Error())
        }
        if gcm {  // GCM
            data, err = symmetric.EncGCM(bytes, pass)
        } else {  // CFB
            data, err = symmetric.EncCFB(bytes, pass)
        }
        if err != nil {
            panic("Error while decrypting data: " + err.Error())
        }
        // Write the unencrypted data
        err = utils.WriteFile(oFilename, data, 0600, b64) // rw for user only
    }

}
