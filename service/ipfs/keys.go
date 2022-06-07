package ipfs

const (
	PROTOCOL        = "https"
	DOMAIN          = "api.pinata.cloud"
	INDEX_ADDR      = PROTOCOL + "://" + DOMAIN
	PINATA_PIN_FILE = INDEX_ADDR + "/pinning/pinFileToIPFS"
	PINATA_PIN_JSON = INDEX_ADDR + "/pinning/pinJSONToIPFS"
)
