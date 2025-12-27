package models

type Person struct {
	PersonID   uint   `gorm:"primaryKey" json:"person_id"`
	PersonName string `json:"person_name"`
	Alias      string `json:"alias"`

	// Relationships
	Transactions []Transaction `gorm:"foreignKey:PersonID" json:"transactions,omitempty"`
}

func (Person) TableName() string {
	return "persons"
}
