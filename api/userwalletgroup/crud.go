package userwalletgroup

import (
	"fmt"
	"moneyplanner/database"
	"moneyplanner/models"
)


// ListWalletGroupsForUser returns all wallet groups reachable from user->wallets->wallet_groups
// Pure ORM approach: Preload relationships and dedupe in Go.
func ListWalletGroupsForUser(userID uint) ([]models.WalletGroup, error) {
	var user models.User

	// Load: user -> wallets -> wallet_groups
	if err := database.DB.
		Preload("Wallets.WalletGroups").
		First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Dedupe wallet groups across multiple wallets
	seen := make(map[uint]models.WalletGroup)

	for _, w := range user.Wallets {
		for _, g := range w.WalletGroups {
			seen[g.WalletGroupID] = g
		}
	}

	groups := make([]models.WalletGroup, 0, len(seen))
	for _, g := range seen {
		groups = append(groups, g)
	}

	return groups, nil
}


// // ListWalletGroupsForUser returns all wallet groups reachable from user->wallets->wallet_groups
// func ListWalletGroupsForUser(userID uint) ([]models.WalletGroup, error) {
// 	var groups []models.WalletGroup

// 	// Join chain:
// 	// users -> user_wallets -> wallets -> wallet_wallet_groups -> wallet_groups
// 	err := database.DB.
// 		Table("wallet_groups").
// 		Select("DISTINCT wallet_groups.*").
// 		Joins("JOIN wallet_wallet_groups ON wallet_wallet_groups.wallet_group_id = wallet_groups.wallet_group_id").
// 		Joins("JOIN wallets ON wallets.wallet_id = wallet_wallet_groups.wallet_id").
// 		Joins("JOIN user_wallets ON user_wallets.wallet_id = wallets.wallet_id").
// 		Where("user_wallets.user_id = ?", userID).
// 		Scan(&groups).Error

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list wallet groups for user: %w", err)
// 	}

// 	return groups, nil
// }
