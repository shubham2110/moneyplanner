package users

import (
	"log"
	"moneyplanner/database"
	"moneyplanner/models"
)

// CreateRootCategories creates the default Income and Expense root categories
func CreateRootCategories(walletID uint) ([]models.Category, error) {
	var rootCategories []models.Category

	incomeCategory := &models.Category{
		Icon:     "💵",
		Name:     "Income",
		WalletID: walletID,
		IsGlobal: true,
		RootID:   0, // Will be set after creation
	}
	if err := database.DB.Create(incomeCategory).Error; err != nil {
		return nil, err
	}
	incomeCategory.RootID = incomeCategory.CategoryID
	database.DB.Save(incomeCategory)
	log.Println("✓ Income root category created")
	rootCategories = append(rootCategories, *incomeCategory)

	expenseCategory := &models.Category{
		Icon:     "💸",
		Name:     "Expense",
		WalletID: walletID,
		IsGlobal: true,
		RootID:   0,
	}
	if err := database.DB.Create(expenseCategory).Error; err != nil {
		return nil, err
	}
	expenseCategory.RootID = expenseCategory.CategoryID
	database.DB.Save(expenseCategory)
	log.Println("✓ Expense root category created")
	rootCategories = append(rootCategories, *expenseCategory)

	return rootCategories, nil
}

// CreateIncomeSubcategories creates default income subcategories
func CreateIncomeSubcategories(walletID uint, parentCategoryID uint) error {
	incomeSubcategories := []string{"Salary", "Refund", "Bonus", "Interest"}
	for _, catName := range incomeSubcategories {
		subcat := &models.Category{
			Name:     catName,
			Icon:     "📊",
			ParentID: &parentCategoryID,
			RootID:   parentCategoryID,
			WalletID: walletID,
			IsGlobal: false,
		}
		if err := database.DB.Create(subcat).Error; err != nil {
			log.Printf("Warning: Failed to create %s subcategory: %v", catName, err)
		} else {
			log.Printf("✓ %s subcategory created", catName)
		}
	}
	return nil
}

// CreateExpenseSubcategories creates default expense subcategories
func CreateExpenseSubcategories(walletID uint, parentCategoryID uint) error {
	expenseSubcategories := []string{"Groceries", "House Maintenance", "Investment", "Utilities", "Transport", "Entertainment"}
	for _, catName := range expenseSubcategories {
		subcat := &models.Category{
			Name:     catName,
			Icon:     "📉",
			ParentID: &parentCategoryID,
			RootID:   parentCategoryID,
			WalletID: walletID,
			IsGlobal: false,
		}
		if err := database.DB.Create(subcat).Error; err != nil {
			log.Printf("Warning: Failed to create %s subcategory: %v", catName, err)
		} else {
			log.Printf("✓ %s subcategory created", catName)
		}
	}
	return nil
}

// CreateAllSubcategories creates both income and expense subcategories
func CreateAllSubcategories(walletID uint, incomeCategoryID uint, expenseCategoryID uint) error {
	if err := CreateIncomeSubcategories(walletID, incomeCategoryID); err != nil {
		return err
	}
	if err := CreateExpenseSubcategories(walletID, expenseCategoryID); err != nil {
		return err
	}
	return nil
}
