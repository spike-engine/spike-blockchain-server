package nft

type nft struct {
	TokenAddress string `json:"token_address"`
	TokenId      string `json:"token_id"`
	Amount       string `json:"amount"`
	ContractType string `json:"contract_type"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	TokenUri     string `json:"token_uri"`
	Metadata     string `json:"metadata"`
}

type FindAllNftsResponse struct {
	Total  int   `json:"total"`
	Result []nft `json:"result"`
}
