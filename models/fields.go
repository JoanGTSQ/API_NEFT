package models

type Field struct {
	NeftModel

	Name string `gorm:"not null"`
}
