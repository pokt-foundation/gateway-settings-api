package models

type Application struct {
	Id              string          `json:"id" bson:"_id"`
	GatewaySettings GatewaySettings `json:"gateway_settings" bson:"gatewaySettings"`
}

type GatewaySettings struct {
	ContractsAllowlist []BlockchainContractsAllowlist `json:"contracts_allowlist" bson:"whitelistContracts"`
}

type BlockchainContractsAllowlist struct {
	BlockchainID string   `json:"blockchain_id" bson:"blockchain_id"`
	Contracts    []string `json:"contracts" bson:"contracts"`
}
