package ipfs

const (
	PROTOCOL        = "https"
	DOMAIN          = "api.pinata.cloud"
	OASISGATEWAY    = "oasis.mypinata.cloud/ipfs/"
	INDEX_ADDR      = PROTOCOL + "://" + DOMAIN
	DOWNLOAD_ADDR   = PROTOCOL + "://" + OASISGATEWAY
	PINATA_PIN_FILE = INDEX_ADDR + "/pinning/pinFileToIPFS"
	PINATA_PIN_JSON = INDEX_ADDR + "/pinning/pinJSONToIPFS"
)
