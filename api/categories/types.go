package categories

import "moneyplanner/models"

type CategoryCreationRequest struct {
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	ParentID *uint  `json:"parent_id,omitempty"`
	WalletID uint   `json:"wallet_id"`
	IsGlobal *bool  `json:"is_global,omitempty"`
}

type CategoryUpdateRequest struct {
	Name     *string `json:"name,omitempty"`
	Icon     *string `json:"icon,omitempty"`
	ParentID *uint   `json:"parent_id,omitempty"`
	IsGlobal *bool   `json:"is_global,omitempty"`
}

// CategoryWithChildren represents a category with its subcategories
type CategoryWithChildren struct {
	Category models.Category        `json:"category"`
	Children []CategoryWithChildren `json:"children"`
}

type CategoryTreeResponse struct {
	Roots []CategoryWithChildren `json:"roots"`
}
