package models

type Target struct {
	NeftModel

	Name string `gorm:"not null"`

	Photo string `gorm:""`

	MissionID int     `gorm:"not null"`
	Mission   Mission `gorm:"-"`
}
