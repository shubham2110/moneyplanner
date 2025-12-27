package users

import (
	"log"
	"moneyplanner/database"
	"moneyplanner/models"
	"time"
)

// CreateWalletGroup creates a new wallet group
func CreateWalletGroup(name string) (*models.WalletGroup, error) {
	walletGroup := &models.WalletGroup{
		WalletGroupName: name,
	}
	if err := database.DB.Create(walletGroup).Error; err != nil {
		return nil, err
	}
	log.Printf("✓ Wallet group '%s' created", name)
	return walletGroup, nil
}

// CreateWallet creates a new wallet with given parameters
func CreateWallet(name, icon string, initialBalance float64) (*models.Wallet, error) {
	wallet := &models.Wallet{
		Name:             name,
		Icon:             icon,
		IsEnabled:        true,
		Balance:          initialBalance,
		LastModifiedTime: time.Now(),
	}
	if err := database.DB.Create(wallet).Error; err != nil {
		return nil, err
	}
	log.Printf("✓ Wallet '%s' created", name)
	return wallet, nil
}

// AssociateWalletWithGroup associates a wallet with a wallet group
func AssociateWalletWithGroup(wallet *models.Wallet, walletGroup *models.WalletGroup) error {
	database.DB.Model(wallet).Association("WalletGroups").Append(walletGroup)
	return nil
}

// AssociateUserWithWallet associates a user with a wallet
func AssociateUserWithWallet(user *models.User, wallet *models.Wallet) error {
	database.DB.Model(user).Association("Wallets").Append(wallet)
	return nil
}
