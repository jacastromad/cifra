# cifra: A simple symmetric encryption tool

`cifra` is an easy-to-use symmetric encryption command-line tool, using the robust AES-256 encryption standard.

Available modes: Galois/Counter Mode (GCM) and Cipher Feedback Mode (CFB).

## Quick Start

```console
$ go build
$ ./cifra -help
```

## Examples

Encrypt a file called data.txt using AES-256 GCM:
```console
$ ./cifra data.txt
Password:
$ ls
data.txt  data.txt.cif
```

Encrypt data.txt using AES-256 CFB. Output file should be named data.txt.cfb:
```console
$ ./cifra -cfb -o data.txt.cfb data.txt
Password:
$ ls
data.txt  data.txt.cif  data.txt.cfb
```

Decrypt both files
```console
$ ./cifra -dec -o data.dec1 data.txt.cif
Password:
$ ./cifra -dec -cfb -o data.dec2 data.txt.cfb
Password:
$ ls
data.txt  data.txt.cif  data.txt.cfb  data.dec1  data.dec2
```

