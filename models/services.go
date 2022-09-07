package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

func NewServices(connectionInfo string, logMode bool) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	db.LogMode(logMode)

	return &Services{
		User:  NewUserService(db),
		Rol:   NewRolService(db),
		Team:  NewTeamService(db),
		Field: NewFieldService(db),
    Device: NewDeviceService(db),
		db:    db,
	}, nil
}

type Services struct {
	User  UserService
	Rol   RolService
	Team  TeamService
	Field FieldService
  Device DeviceService
	db    *gorm.DB
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) DestructiveReset() error {
	if err := s.db.DropTableIfExists(&Rol{}, &pwReset{}, &User{},
		&Team{}, &Division{}, &Target{},
		&Mission{},
		&Field{}, &War{}, &Device{}).Error; err != nil {
		return err
	}
	return s.AutoMigrate()
}

func (s *Services) DestructiveStatic() error {
	if err := s.db.DropTableIfExists().Error; err != nil {
		return err
	}
	return s.AutoMigrate()
}

func (s *Services) AutoMigrate() error {
	if err := s.db.AutoMigrate(&User{}, &pwReset{}, &Rol{},
		&Field{}, &War{},
		&Mission{}, &Target{}, &Division{},
		&Team{}, &Device{}).Error; err != nil {
		return err
	}
	return nil
}

type NeftModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}
