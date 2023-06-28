package main

import "fmt"
import "encoding/base64"
import "crypto/aes"
import "crypto/cipher"
import "os"
import "path/filepath"
import "strings"
import "time"

var bytes = []byte{19, 64, 55, 23, 99, 51, 87, 34, 13, 12, 79, 72, 41, 8, 3, 81}

func encrypt( message, key string ) string {
    
    // implementing code from https://blog.logrocket.com/learn-golang-encryption-decryption/

    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        panic(err)
    }
    
    plainText := []byte(message)
    cfb := cipher.NewCFBEncrypter(block, bytes)
    cipherText := make([]byte, len(plainText))
    cfb.XORKeyStream(cipherText, plainText)
    
    return base64.StdEncoding.EncodeToString(cipherText)
}

func main() {
    // grab the filename from the command-line argument
    if len(os.Args[1:]) == 0 {
        fmt.Printf("Usage %s <filename> [filename2]   ; compiles filename into filename.o.txt\n", filepath.Base(os.Args[0]))
        return
    }
    
    for i := 0; i < len(os.Args[1:]); i++ {
        // start of path-checking code, but for now just clean for ease-of-visibility
        filename := filepath.Clean(os.Args[1+i])
        
        // read in the code
        dat, err := os.ReadFile(filename)
        if err != nil {
            fmt.Printf("Error [%s]> %s\n", filename, err.Error())
            continue
        }
        
        // compile the code
        compiled := encrypt(string(dat), "Compiler_v1.0*@34Key!!!$")
        
        // create the program name and timestamp header, for convenience
        // <filename> <timestamp>   is the header ,
        // <code goes here>
        curtime := time.Now()
        header := fmt.Sprintf("%s %d:%d %d %s %d\n", filepath.Base(filename), curtime.Hour(), curtime.Minute(), curtime.Day(), curtime.Month(), curtime.Year())
        
        // determine the output filename
        outpathfile := strings.TrimSuffix(filename, filepath.Ext(filename))
        
        
        // add regular spacing to the compiled file for ease-of-printing
        var formatted string
        linelength := 50
        textpreform := compiled
        linecount := (len(textpreform) / linelength) + 1
        
        for i := 0; i < linecount; i++ {
            if i*linelength >= len(textpreform) {
                break
            }
            
            remainder := len(textpreform) - (i*linelength)
            if remainder > linelength {
                remainder = linelength
            }
            
            line := textpreform[i*linelength:i*linelength+remainder] + "\n"
            
            formatted = formatted + line
        }
        
        // trim off the last new-line
        formatted = strings.TrimSuffix(formatted, "\n")
        
        // add the header after the encrypt-compiled text, since it already contains a newline
        formatted = header + formatted

        
        // write the encoded form to file
        err = os.WriteFile(outpathfile + ".o.txt", []byte(formatted), 0644)
        if err != nil {
            fmt.Printf("Error [%s]> %s\n", outpathfile, err.Error())
            continue
        }
        
        // print status message, before moving on to next file!
        fmt.Printf("Successfully compiled %s\n", filename)
    }
}