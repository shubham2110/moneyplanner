package wallet

type WalletCreationRequest struct {
	Name      string   `json:"name"`
	Icon      string   `json:"icon"`
	IsEnabled *bool    `json:"is_enabled,omitempty"`
	Balance   *float64 `json:"balance,omitempty"`
}

type WalletUpdateRequest struct {
	Name      *string  `json:"name,omitempty"`
	Icon      *string  `json:"icon,omitempty"`
	IsEnabled *bool    `json:"is_enabled,omitempty"`
	Balance   *float64 `json:"balance,omitempty"`
}
