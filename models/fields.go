package models

import (
	"regexp"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"neft.web/hash"
)

type FieldDB interface {
	Create(field *Field) error
	Update(field *Field) error
	Delete(field *Field) error
}

type FieldService interface {
	FieldDB
}

func newFieldGorm(db *gorm.DB) (*fieldGorm, error) {
	return &fieldGorm{
		db: db,
	}, nil
}
func NewFieldService(gD *gorm.DB) FieldService {
	ug, err := newFieldGorm(gD)
	if err != nil {
		return nil
	}
	hmac := hash.NewHMAC(hmacScretKey)
	uv := newFieldValidator(ug, hmac)
	return &fieldService{
		FieldDB: uv,
	}
}

type fieldService struct {
	FieldDB
}
type fieldValidator struct {
	FieldDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

func newFieldValidator(udb FieldDB, hmac hash.HMAC) *fieldValidator {
	return &fieldValidator{
		FieldDB:    udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

var _ FieldDB = &fieldGorm{}

type fieldGorm struct {
	db *gorm.DB
}

func (tg *fieldGorm) Create(field *Field) error {
	err := tg.db.Create(field).Error
	if err != nil {
		return err
	}
	return nil
}

func (tg *fieldGorm) Delete(field *Field) error {
	return tg.db.Delete(&field).Error
}

func (tg *fieldGorm) Update(field *Field) error {
	return tg.db.Save(field).Error
}



type Field struct {
	NeftModel

	Name   string `gorm:"not null"`
	Coords string
}
