## How to start spike-blockchain-service?

#### 1. Create a file named .env and Configure these options

```
MORALIS_KEY="7L9fFRFaUpKOlFoUZi7SN0UBiUWbHgf6tTenyqg19TjPRTGuRS46UMRierHQ3XIs"

PINATA_API_KEY="6d2a1a3f3775a79ab198"
PINATA_SECRET_KEY="5224d1e3a6388fe9ca32f3b7b49e3877ac02a01049a1f0b0c62944fe7aa0f542"
PINATA_GATEWAY="oasis.mypinata.cloud"
```

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

