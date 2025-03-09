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
-t (int) -- number of threads
-p (str) -- path to csv file with wallets
```

You can use [Clickhouse](https://github.com/ClickHouse/ClickHouse) for storing generated wallets:
```
apt-get -y install apt-transport-https ca-certificates curl gnupg
curl -fsSL 'https://packages.clickhouse.com/rpm/lts/repodata/repomd.xml.key' | sudo gpg --dearmor -o /usr/share/keyrings/clickhouse-keyring.gpg
ARCH=$(dpkg --print-architecture)
echo "deb [signed-by=/usr/share/keyrings/clickhouse-keyring.gpg arch=${ARCH}] https://packages.clickhouse.com/deb stable main" | sudo tee /etc/apt/sources.list.d/clickhouse.list
apt-get update
apt-get -y install clickhouse-server clickhouse-client
```

Create ReplacingMergeTree table to avoid duplicates:
```
CREATE TABLE default.kaspa
(
    `dt` DateTime DEFAULT now(),
    `address` String,
    `private` String
)
ENGINE = ReplacingMergeTree
PARTITION BY dt
ORDER BY address
SETTINGS index_granularity = 8192
```

Export csv file to the default.kaspa Clickhouse table:
```
cat /root/kaspa.csv | clickhouse-client --host clickhouse --port 9000 --query='insert into default.kaspa (address,private) format CSV'
```

Select from Clickhouse:
```
clickhouse-client
select count(address) from kaspa
select address,private from kaspa limit 100 format csv
select address,private from kaspa where address ilike 'kaspa:qqqqq%' limit 10
```

### Always check that generated seed-phrase corresponds to the specified address
