package main

import (
    "fmt"
    "log"
    "flag"
    "os"
    "path/filepath"

    "math/rand"
    "time"

    "sync"

    "github.com/kaspanet/kaspad/cmd/kaspawallet/libkaspawallet"
    "github.com/kaspanet/kaspad/cmd/kaspawallet/keys"
    "github.com/kaspanet/kaspad/domain/dagconfig"
)

func generatePassword() string {
    str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    bytes := []byte(str)
    result := []byte{}
    rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
    for i := 0; i < 16; i++ {
        result = append(result, bytes[rand.Intn(len(bytes))])
    }
return string(result)
}

func generateWallet( wg *sync.WaitGroup, dbpath string ) {
    defer wg.Done()
    
    for {
        password := generatePassword()
        
        // create bip-39 mnemonic
        mnemonic, err := libkaspawallet.CreateMnemonic()
        if err != nil {
            log.Fatalf("mnemonic generation error: %v", err)
        }
        
        // generate key pair
        keysFile, err := keys.NewFileFromMnemonic(&dagconfig.MainnetParams, mnemonic, password)
        if err != nil {
            log.Fatalf("keysFile error: %v", err)
        }
        
        // extract first address
        path := fmt.Sprintf("m/%d/%d", libkaspawallet.ExternalKeychain, 0)
        address, err := libkaspawallet.Address(&dagconfig.MainnetParams, keysFile.ExtendedPublicKeys, 1, path, false)
        if err != nil {
            log.Fatalf("address generation error: %v", err)
        }
        
        // print the result in csv format
        // fmt.Printf("%v,%v\n", address,mnemonic)

        // write to file instead
        csvLine := fmt.Sprintf("%v,%v\n", address, mnemonic)
        f, _ := os.OpenFile(dbpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
        f.Write([]byte(csvLine))
        f.Close()
    }
}

func main() { 
    numPtr := flag.Int("t", 1, "-t must be an integer")

    homeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Fprintf(os.Stderr, "error while getting home path: %v\n", err)
        os.Exit(1)
    }
    defaultPath := filepath.Join(homeDir, "kaspa.csv")
    filePath := flag.String("p", defaultPath, "path to kaspa.csv file")

    flag.Parse()
    
    num := *numPtr
    dbpath := *filePath
    
    var wg sync.WaitGroup
    
    wg.Add(num)
    
    for i := 0; i < num; i++ {
        go generateWallet( &wg, dbpath )
    }
    
    wg.Wait()
}
