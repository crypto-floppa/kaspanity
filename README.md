### Always check that generated seed-phrase corresponds to the specified address

Install git and golang packages.
Run `go install` for building the binary.
```
apt-get -y install git golang
git clone https://github.com/crypto-floppa/kaspanity.git
cd kaspanity
go install
```

kaspanity uses two options, you can change default values by flags: `/root/go/bin/kaspanity -t 5 -p /root/kaspa.csv`, where:
```
-w (int) -- number of wallets
-a (int) -- number of addresses
```
### Always check that generated seed-phrase corresponds to the specified address
