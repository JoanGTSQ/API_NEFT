package models

import (
	"regexp"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"neft.web/hash"
)

type FieldDB interface {
	ByID(field *Field) error
	AllFields() ([]*Field, error)

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

func (tg *fieldGorm) ByID(field *Field) error {

	db := tg.db.Where("id = ?", field.ID).
		First(&field).Error
	return db
}

func (tg *fieldGorm) AllFields() ([]*Field, error) {
	var fields []*Field
	err := tg.db.Find(&fields).Error
	if err != nil {
		return nil, err
	}
	return fields, nil
}

type Field struct {
	NeftModel

	Name   string `gorm:"not null" json:"name"`
	Coords string `gorm:"not null" json:"coords"`
}
