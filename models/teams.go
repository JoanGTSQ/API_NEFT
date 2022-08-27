package models

import (
	"regexp"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"neft.web/hash"
)

type TeamDB interface {
	ByID(id uint) (*Team, error)
	AllTeamByID(id uint) (*Team, error)

	Create(team *Team) error
	Update(team *Team) error
	Delete(id string) error
}

type TeamService interface {
	TeamDB
}

func newTeamGorm(db *gorm.DB) (*teamGorm, error) {
	return &teamGorm{
		db: db,
	}, nil
}
func NewTeamService(gD *gorm.DB) TeamService {
	ug, err := newTeamGorm(gD)
	if err != nil {
		return nil
	}
	hmac := hash.NewHMAC(hmacScretKey)
	uv := newTeamValidator(ug, hmac)
	return &teamService{
		TeamDB: uv,
	}
}

type teamService struct {
	TeamDB
}
type teamValidator struct {
	TeamDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

func newTeamValidator(udb TeamDB, hmac hash.HMAC) *teamValidator {
	return &teamValidator{
		TeamDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

var _ TeamDB = &teamGorm{}

type teamGorm struct {
	db *gorm.DB
}

func (tg *teamGorm) Create(team *Team) error {
	err := tg.db.Create(team).Error
	if err != nil {
		return err
	}
	return nil
}

func (tg *teamGorm) Delete(id string) error {
	varInt, err := strconv.Atoi(id)
	if err != nil {
		return ERR_ID_INVALID
	}
	team := Team{NeftModel: NeftModel{ID: uint(varInt)}}
	return tg.db.Delete(&team).Error
}

func (tg *teamGorm) Update(team *Team) error {
	return tg.db.Save(team).Error
}

func (ug *teamGorm) ByID(id uint) (*Team, error) {
	var team Team
	db := ug.db.Where("id = ?", id).First(&team)
	err := first(db, &team)
	return &team, err
}

// SEARCH BY ID
func (ug *teamGorm) AllTeamByID(id uint) (*Team, error) {
	var team Team
	db := ug.db.Where("id = ?", id).
		Preload("TeamLead").
		// Preload("TeamLead.Rol").
		Preload("Member1").
		// Preload("Member1.Rol").
		Preload("Member2").
		// Preload("Member2.Rol").
		Preload("Member3").
		// Preload("Member3.Rol").
		Preload("Member4").
		// Preload("Member4.Rol").
		Preload("Member5").
		// Preload("Member5.Rol").
		Preload("Division1").
		Preload("Division1.Commandant").
		Preload("AssignedMission").
		First(&team).Error
	return &team, db
}

type Team struct {
	NeftModel

	Name string `gorm:"not null"`

	TeamLeadID int  `gorm:"unique_index" json:"teamleadid"`
	TeamLead   User `gorm:"foreignkey:TeamLeadId" json:"teamlead"`

	Member1ID int  `gorm:"unique_index" json:"member1id"`
	Member1   User `gorm:"foreignkey:Member1ID" json:"member1"`

	Member2ID int  `gorm:"unique_index" json:"member2id"`
	Member2   User `gorm:"foreignkey:Member2ID" json:"member2"`

	Member3ID int  `gorm:"unique_index" json:"member3id"`
	Member3   User `gorm:"foreignkey:Member3ID" json:"member3"`

	Member4ID int  `gorm:"unique_index" json:"member4id"`
	Member4   User `gorm:"foreignkey:Member4ID" json:"member4"`

	Member5ID int  `gorm:"unique_index" json:"member5id"`
	Member5   User `gorm:"foreignkey:Member5ID" json:"member5"`

	Division1ID int      `gorm:"" json:"division1id"`
	Division1   Division `gorm:"foreignkey:division1ID" json:"division"`

	AssignedMissionID int     `gorm:"" json:"assignedmissionid"`
	AssignedMission   Mission `gorm:"foreignkey:AssignedMissionID" json:"assignedmission"`
}
