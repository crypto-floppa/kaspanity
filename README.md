### Always check that generated seed-phrase corresponds to the specified address

Install git and golang packages.
Run `go install` for building the binary.
```
apt-get -y install git golang
git clone https://github.com/crypto-floppa/kaspanity.git
cd kaspanity
go install
```

kaspanity uses two options, you can change default values by flags: `/root/go/bin/kaspanity -w 10 -a 10000`, where:
```
-w (int) -- number of wallets   (default is 1)
-a (int) -- number of addresses (default is 2)
```

### Always check that generated seed-phrase corresponds to the specified address
