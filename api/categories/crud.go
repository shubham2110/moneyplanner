package categories

import (
	"fmt"
	"log"
	"moneyplanner/database"
	"moneyplanner/models"
)

// GetCategoryByID retrieves a category by its ID
func GetCategoryByID(categoryID uint) (*models.Category, error) {
	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	// Load wallet
	var wallet models.Wallet
	if err := database.DB.First(&wallet, category.WalletID).Error; err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	category.Wallet = wallet
	return &category, nil
}

// ListCategoriesByWallet retrieves all categories for a wallet (flat list)
func ListCategoriesByWallet(walletID uint) ([]models.Category, error) {
	var categories []models.Category
	if err := database.DB.Where("wallet_id = ?", walletID).Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	// Load wallet
	var wallet models.Wallet
	if err := database.DB.First(&wallet, walletID).Error; err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	// Set wallet in each category
	for i := range categories {
		categories[i].Wallet = wallet
	}
	return categories, nil
}

// GetCategoryTree retrieves categories in hierarchical structure for a wallet
func GetCategoryTree(walletID uint) (*CategoryTreeResponse, error) {
	// Load wallet
	var wallet models.Wallet
	if err := database.DB.First(&wallet, walletID).Error; err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Get root categories (ParentID is NULL)
	var roots []models.Category
	if err := database.DB.Where("wallet_id = ? AND parent_id IS NULL", walletID).Find(&roots).Error; err != nil {
		return nil, fmt.Errorf("failed to get root categories: %w", err)
	}

	var rootWithChildren []CategoryWithChildren
	for _, root := range roots {
		root.Wallet = wallet // Set wallet
		children, err := buildChildren(root.CategoryID, wallet)
		if err != nil {
			return nil, err
		}
		rootWithChildren = append(rootWithChildren, CategoryWithChildren{
			Category: root,
			Children: children,
		})
	}

	return &CategoryTreeResponse{Roots: rootWithChildren}, nil
}

// buildChildren recursively builds the children for a category
func buildChildren(parentID uint, wallet models.Wallet) ([]CategoryWithChildren, error) {
	var children []models.Category
	if err := database.DB.Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
		return nil, fmt.Errorf("failed to get children for category %d: %w", parentID, err)
	}

	var result []CategoryWithChildren
	for _, child := range children {
		child.Wallet = wallet // Set wallet
		grandChildren, err := buildChildren(child.CategoryID, wallet)
		if err != nil {
			return nil, err
		}
		result = append(result, CategoryWithChildren{
			Category: child,
			Children: grandChildren,
		})
	}

	return result, nil
}

// CreateCategory creates a new category
func CreateCategory(req *CategoryCreationRequest) (*models.Category, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.WalletID == 0 {
		return nil, fmt.Errorf("wallet_id is required")
	}

	c := &models.Category{
		Name:     req.Name,
		Icon:     req.Icon,
		ParentID: req.ParentID,
		WalletID: req.WalletID,
		IsGlobal: false, // default
	}

	// Set defaults
	if req.IsGlobal != nil {
		c.IsGlobal = *req.IsGlobal
	}

	// Create category
	if err := database.DB.Create(c).Error; err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	// Set RootID
	if req.ParentID == nil {
		// This is a root category
		c.RootID = c.CategoryID
	} else {
		// Get parent's RootID
		parent, err := GetCategoryByID(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent category: %w", err)
		}
		c.RootID = parent.RootID
	}

	// Update RootID
	if err := database.DB.Model(c).Update("root_id", c.RootID).Error; err != nil {
		return nil, fmt.Errorf("failed to update root_id: %w", err)
	}

	// If global, sync to all wallets
	if c.IsGlobal {
		if err := SyncGlobalCategoryToAllWallets(c.CategoryID); err != nil {
			log.Printf("Warning: Failed to sync global category: %v", err)
		}
	}

	log.Printf("✓ Category '%s' created (ID: %d)", c.Name, c.CategoryID)
	return c, nil
}

// UpdateCategory updates category details
func UpdateCategory(categoryID uint, req *CategoryUpdateRequest) (*models.Category, error) {
	category, err := GetCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}

	if req.Name != nil {
		updates["name"] = *req.Name
		category.Name = *req.Name
	}

	if req.Icon != nil {
		updates["icon"] = *req.Icon
		category.Icon = *req.Icon
	}

	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
		category.ParentID = req.ParentID
		// Update RootID
		if *req.ParentID == 0 {
			// Becoming root
			updates["root_id"] = categoryID
			category.RootID = categoryID
		} else {
			parent, err := GetCategoryByID(*req.ParentID)
			if err != nil {
				return nil, fmt.Errorf("failed to get parent category: %w", err)
			}
			updates["root_id"] = parent.RootID
			category.RootID = parent.RootID
		}
	}

	if req.IsGlobal != nil {
		updates["is_global"] = *req.IsGlobal
		category.IsGlobal = *req.IsGlobal
	}

	if len(updates) == 0 {
		return category, nil // No updates provided
	}

	if err := database.DB.Model(&models.Category{}).
		Where("category_id = ?", categoryID).
		Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	// If set to global, sync to all wallets
	if isGlobalUpdated, ok := updates["is_global"]; ok && isGlobalUpdated == true {
		if err := SyncGlobalCategoryToAllWallets(categoryID); err != nil {
			log.Printf("Warning: Failed to sync global category: %v", err)
		}
	}

	log.Printf("✓ Category '%s' (ID: %d) updated", category.Name, categoryID)
	return category, nil
}

// DeleteCategory deletes a category by ID
func DeleteCategory(categoryID uint) error {
	// Check if category exists
	category, err := GetCategoryByID(categoryID)
	if err != nil {
		return err
	}

	// Check if has children
	var childrenCount int64
	if err := database.DB.Model(&models.Category{}).Where("parent_id = ?", categoryID).Count(&childrenCount).Error; err != nil {
		return fmt.Errorf("failed to check children: %w", err)
	}
	if childrenCount > 0 {
		return fmt.Errorf("cannot delete category with children")
	}

	// Delete the category
	if err := database.DB.Delete(&models.Category{}, categoryID).Error; err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	log.Printf("✓ Category '%s' (ID: %d) deleted", category.Name, categoryID)
	return nil
}

// SyncGlobalCategoryToAllWallets syncs a global category to all wallets that don't have it
func SyncGlobalCategoryToAllWallets(categoryID uint) error {
	category, err := GetCategoryByID(categoryID)
	if err != nil {
		return err
	}

	if !category.IsGlobal {
		return fmt.Errorf("category is not global")
	}

	// Get all wallets
	var wallets []models.Wallet
	if err := database.DB.Find(&wallets).Error; err != nil {
		return fmt.Errorf("failed to get wallets: %w", err)
	}

	for _, wallet := range wallets {
		// Check if category exists for this wallet
		var count int64
		database.DB.Model(&models.Category{}).Where("wallet_id = ? AND name = ? AND is_global = ?", wallet.WalletID, category.Name, true).Count(&count)
		if count == 0 {
			// Create
			newCat := &models.Category{
				Name:     category.Name,
				Icon:     category.Icon,
				WalletID: wallet.WalletID,
				IsGlobal: true,
				RootID:   0, // Will be set after create
			}
			if err := database.DB.Create(newCat).Error; err != nil {
				log.Printf("Warning: Failed to create global category for wallet %d: %v", wallet.WalletID, err)
			} else {
				newCat.RootID = newCat.CategoryID
				database.DB.Save(newCat)
				log.Printf("✓ Global category '%s' synced to wallet %d", category.Name, wallet.WalletID)
			}
		}
	}

	return nil
}
