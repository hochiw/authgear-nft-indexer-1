package model

type RawContract struct {
	Value   string `json:"value"`
	Address string `json:"address"`
	Decimal string `json:"decimal"`
}

type Metadata struct {
	BlockTimestamp string `json:"blockTimestamp"`
}

type TokenTranfer struct {
	Category        string      `json:"category"`
	Token           string      `json:"token"`
	BlockNum        string      `json:"blockNum"`
	From            string      `json:"from"`
	To              string      `json:"to"`
	Value           string      `json:"value"`
	ERC721TokenID   string      `json:"erc721TokenId"`
	ERC1155Metadata string      `json:"erc1155Metadata"`
	TokenID         string      `json:"tokenId"`
	Asset           string      `json:"asset"`
	Hash            string      `json:"hash"`
	RawContract     RawContract `json:"rawContract"`
	Metadata        Metadata    `json:"metadata"`
}

type AssetTransferRequestParams struct {
	FromBlock         string   `json:"fromBlock,omitempty"`
	ToBlock           string   `json:"toBlock,omitempty"`
	FromAddress       string   `json:"fromAddress,omitempty"`
	ToAddress         string   `json:"toAddress,omitempty"`
	ContractAddresses []string `json:"contractAddresses,omitempty"`
	Category          []string `json:"category"`
	WithMetadata      bool     `json:"withMetadata,omitempty"`
	ExcludeZeroValue  bool     `json:"excludeZeroValue,omitempty"`
	MaxCount          string   `json:"maxCount,omitempty"`
	PageKey           string   `json:"pageKey,omitempty"`
}

type AssetTransferResult struct {
	Transfers []TokenTranfer `json:"transfers"`
	PageKey   string         `json:"pageKey"`
}

type AssetTransferError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AssetTransferResponse struct {
	Result AssetTransferResult `json:"result"`
	Error  *AssetTransferError `json:"error,omitempty"`
}
