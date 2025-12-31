package persons

import (
	"fmt"
	"log"
	"moneyplanner/database"
	"moneyplanner/models"
)

// GetPersonByID retrieves a person by its ID
func GetPersonByID(personID uint) (*models.Person, error) {
	var person models.Person
	if err := database.DB.First(&person, personID).Error; err != nil {
		return nil, fmt.Errorf("person not found: %w", err)
	}
	return &person, nil
}

// GetPersonByName retrieves a person by name
func GetPersonByName(personName string) (*models.Person, error) {
	var person models.Person
	if err := database.DB.Where("person_name = ?", personName).First(&person).Error; err != nil {
		return nil, fmt.Errorf("person not found: %w", err)
	}
	return &person, nil
}

// ListAllPersons retrieves all persons
func ListAllPersons() ([]models.Person, error) {
	var persons []models.Person
	if err := database.DB.Find(&persons).Error; err != nil {
		return nil, fmt.Errorf("failed to list persons: %w", err)
	}
	return persons, nil
}

// CreatePerson creates a new person
func CreatePerson(req *PersonCreationRequest) (*models.Person, error) {
	if req.PersonName == "" {
		return nil, fmt.Errorf("person_name is required")
	}

	p := &models.Person{
		PersonName: req.PersonName,
		Alias:      req.Alias,
	}

	if err := database.DB.Create(p).Error; err != nil {
		return nil, fmt.Errorf("failed to create person: %w", err)
	}

	log.Printf("✓ Person '%s' created (ID: %d)", p.PersonName, p.PersonID)
	return p, nil
}

// UpdatePerson updates person details
func UpdatePerson(personID uint, req *PersonUpdateRequest) (*models.Person, error) {
	person, err := GetPersonByID(personID)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}

	if req.PersonName != nil {
		updates["person_name"] = *req.PersonName
		person.PersonName = *req.PersonName
	}

	if req.Alias != nil {
		updates["alias"] = *req.Alias
		person.Alias = *req.Alias
	}

	if len(updates) == 0 {
		return person, nil // No updates provided
	}

	if err := database.DB.Model(&models.Person{}).
		Where("person_id = ?", personID).
		Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update person: %w", err)
	}

	log.Printf("✓ Person '%s' (ID: %d) updated", person.PersonName, personID)
	return person, nil
}

// DeletePerson deletes a person by ID
func DeletePerson(personID uint) error {
	// Check if person exists
	person, err := GetPersonByID(personID)
	if err != nil {
		return err
	}

	// Delete the person
	if err := database.DB.Delete(&models.Person{}, personID).Error; err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}

	log.Printf("✓ Person '%s' (ID: %d) deleted", person.PersonName, personID)
	return nil
}
