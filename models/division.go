package models

type Division struct {
	NeftModel

	Name string `gorm:"not null"`

	Flag string `gorm:""`

	CommandantID int  `json:"commandantid"`
	Commandant   User `gorm:"foreignkey:commandantID" json:"commandant"`

	WarID int `json:"warid"`
	War   War `gorm:"foreignkey:warID" json:"war"`
}
