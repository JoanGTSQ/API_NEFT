package models

import (
	"regexp"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"neft.web/hash"
)

type RolDB interface {
	RetrieveAll() ([]*Rol, error)
	ReturnFiveUserByRolID(id string) ([]*User, error)
	ReturnUserRol() ([]*UserRol, error)
}

type RolService interface {
	RolDB
}

func newRolGorm(db *gorm.DB) (*rolGorm, error) {
	return &rolGorm{
		db: db,
	}, nil
}
func NewRolService(gD *gorm.DB) RolService {
	ug, err := newRolGorm(gD)
	if err != nil {
		return nil
	}
	hmac := hash.NewHMAC(hmacScretKey)
	uv := newRolValidator(ug, hmac)
	return &rolService{
		RolDB: uv,
	}
}

type rolService struct {
	RolDB
}
type rolValidator struct {
	RolDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

func newRolValidator(udb RolDB, hmac hash.HMAC) *rolValidator {
	return &rolValidator{
		RolDB:      udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

var _ RolDB = &rolGorm{}

type rolGorm struct {
	db *gorm.DB
}

func (ug *rolGorm) ByID(id uint) (*Rol, error) {
	var rol Rol
	db := ug.db.Where("id = ?", id).First(&rol)
	err := first(db, &rol)
	return &rol, err
}

func (ug *rolGorm) RetrieveAll() ([]*Rol, error) {
	var roles []*Rol
	err := ug.db.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (ug *rolGorm) ReturnFiveUserByRolID(id string) ([]*User, error) {
	user := make([]*User,5)

	err := ug.db.Where("role_id = ?", id).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

type UserRol struct {
	Id    int
	Count int
}

func (ug *rolGorm) ReturnUserRol() ([]*UserRol, error) {

	var userCount []*UserRol
	for i := 1; i <= 6; i++ {
		var countless UserRol
		countless.Id = i
		ug.db.Model(&User{}).Where("role_id = ?", i).Count(&countless.Count)
		userCount = append(userCount, &countless)
	}

	return userCount, nil
}

type Rol struct {
	NeftModel

	Name string `gorm:"not null" json:"name"`

	CreateTeams       bool `gorm:"not null;default:false"`
	DeleteTeams       bool `gorm:"not null;default:false"`
	ManageTeamMembers bool `gorm:"not null;default:false"`

	CreateExternalSoldier bool `gorm:"not null;default:false"`
	DeleteExternalSoldier bool `gorm:"not null;default:false"`
	ModifySoldierData     bool `gorm:"not null;default:false"`

	CreateTarget   bool `gorm:"not null;default:false"`
	DeleteTarget   bool `gorm:"not null;default:false"`
	AssignTarget   bool `gorm:"not null;default:false"`
	ModifyTarget   bool `gorm:"not null;default:false"`
	CompleteTarget bool `gorm:"not null;default:false"`

	CreateField bool `gorm:"not null;default:false"`
	DeleteField bool `gorm:"not null;default:false"`

	CreateMissions   bool `gorm:"not null;default:false"`
	DeleteMissions   bool `gorm:"not null;default:false"`
	ModifyMissions   bool `gorm:"not null;default:false"`
	CompleteMissions bool `gorm:"not null;default:false"`

	CreateWar   bool `gorm:"not null;default:false"`
	DeleteWar   bool `gorm:"not null;default:false"`
	ModifyWar   bool `gorm:"not null;default:false"`
	InviteToWar bool `gorm:"not null;default:false"`
}
