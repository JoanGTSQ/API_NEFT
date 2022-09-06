package models

import (
	"regexp"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"neft.web/hash"
)

type WarDB interface {
	ByID(id uint) (*War, error)
}

type WarService interface {
	WarDB
}

func newWarGorm(db *gorm.DB) (*warGorm, error) {
	return &warGorm{
		db: db,
	}, nil
}
func NewWarService(gD *gorm.DB) WarService {
	ug, err := newWarGorm(gD)
	if err != nil {
		return nil
	}
	hmac := hash.NewHMAC(hmacScretKey)
	uv := newWarValidator(ug, hmac)
	return &warService{
		WarDB: uv,
	}
}

type warService struct {
	WarDB
}
type warValidator struct {
	WarDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

func newWarValidator(udb WarDB, hmac hash.HMAC) *warValidator {
	return &warValidator{
		WarDB:      udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

var _ WarDB = &warGorm{}

type warGorm struct {
	db *gorm.DB
}

func (ug *warGorm) ByID(id uint) (*War, error) {
	var war War
	db := ug.db.Where("id = ?", id).First(&war)
	err := first(db, &war)
	return &war, err
}

type War struct {
	NeftModel

	Name string `gorm:"not null"`

	FieldID int   `json:"fieldid"`
	Field   Field `gorm:"foreignkey:fieldID" json:"field"`

	StartDate time.Time
	EndDate   time.Time
}
