## How to start spike-blockchain-service?

#### 1. Create a file named .env and Configure these options

```
MORALIS_KEY="XXXXXXXXXXXXXXX"

PINATA_API_KEY="XXXXXXXXXXXX"
PINATA_SECRET_KEY="XXXXXXXXX"
PINATA_GATEWAY="XXXXXXXXXXXX"
```

##### 1.1 What is Pinata and how to get PINATA_X?

[Pinata](https://app.pinata.cloud/pinmanager) is a third-party storage service provider of IPFS.

PINATA_API_KEY & PINATA_SECRET_KEY Used for authentication, PINATA_GATEWAY Private gateways can be created. The advantage is that they are fast and can also use public gateways.

```
public gateways :
PINATA_GATEWAY="api.pinata.cloud"
```



PINATA_X needs to be created in the user's console, with which you can start using the services provided by Pinata.



##### 1.2 What is Moralis and how to get MORALIS_KEY

moralis is cross-chain web3 service provider,You can create a key for access in the user console.

#### 2.  Change the startup port number

open spike-blockchain-service/main.go

```
func main() {
	config.Init()

	r := server.NewRouter()
	r.Run(":40000")
}
```

#### 3. Compile the project to use

