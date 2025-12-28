package walletgroup

type WalletGroupCreationRequest struct {
	WalletGroupName string `json:"wallet_group_name"`
}

type WalletGroupUpdateRequest struct {
	WalletGroupName *string `json:"wallet_group_name,omitempty"`
}
