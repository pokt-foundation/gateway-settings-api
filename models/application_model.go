package models

type Application struct {
	Id              string          `json:"id" bson:"_id"`
	GatewaySettings GatewaySettings `json:"gateway_settings" bson:"gateway_settings"`
}

type GatewaySettings struct {
	AllowlistContracts []string `json:"allowlist_contracts,omitempty" validate:"required"`
}
