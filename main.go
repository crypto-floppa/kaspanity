package main

import (
    "fmt"
    "log"
    "flag"

    "math/rand"
    "encoding/json"
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

func generateWallet( wg *sync.WaitGroup, awallet int ) {
    defer wg.Done()

    // for {
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

        // extract multiple addresses
        addresses := make([]string, 0, awallet)
        for i := 0; i < awallet; i++ {
            path := fmt.Sprintf("m/%d/%d", libkaspawallet.ExternalKeychain, i)
            address, err := libkaspawallet.Address(&dagconfig.MainnetParams, keysFile.ExtendedPublicKeys, 1, path, false)
            if err != nil {
                log.Fatalf("address generation error for index %d: %v", i, err)
            }
        addresses = append(addresses, address.String())
        }

        output := map[string]interface{}{
            "mnemonic": mnemonic,
            "addresses": addresses,
        }

        jsonOutput, err := json.MarshalIndent(output, "", "  ")
        if err != nil {
            log.Fatalf("failed to marshal to JSON: %v", err)
        }

        fmt.Println(string(jsonOutput))

    // }
}

func main() {
    numW := flag.Int("w", 1, "-w must be an integer")
    numA := flag.Int("a", 2, "-a must be an integer")

    flag.Parse()

    nwallet := *numW
    awallet := *numA

    var wg sync.WaitGroup

    wg.Add(nwallet)

    for i := 0; i < nwallet; i++ {
        go generateWallet( &wg, awallet )
    }

    wg.Wait()
}
