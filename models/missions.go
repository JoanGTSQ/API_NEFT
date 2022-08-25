package models

type Mission struct {
	NeftModel

	Name string `gorm:"not null"`

	Complete bool `gorm:"default: false"`
}








