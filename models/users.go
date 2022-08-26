package models

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"neft.web/hash"
	"neft.web/rand"
)

const (
	userPwPPepper = "JUNneft"
	hmacScretKey  = "sec2022"
)

type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	GetAllUsers() ([]*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id string) error
}

type UserService interface {
	Authenticate(email, password string) (*User, error)

	InitiateReset(userID string) (string, error)
	CompleteReset(token, newPw string) (*User, error)
	UserDB
}

func NewUserService(gD *gorm.DB) UserService {
	ug, err := newUserGorm(gD)
	if err != nil {
		return nil
	}
	hmac := hash.NewHMAC(hmacScretKey)
	uv := newUserValidator(ug, hmac)
	return &userService{
		UserDB:    uv,
		pwResetDB: newPwResetValidator(&pwResetGorm{db: gD}, hmac),
	}
}

type userService struct {
	UserDB
	pwResetDB pwResetDB
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user, uv.normalizeEmail, uv.defaultify, uv.hmacRemember); err != nil {
		return nil, err
	}

	return uv.UserDB.ByEmail(user.Email)
}

func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}

	pwByte := []byte(user.Password + userPwPPepper)

	hashedBytes, err := bcrypt.GenerateFromPassword(pwByte, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ERR_PSSWD_TOO_SHORT
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ERR_PSSWD_REQUIRED
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ERR_PSSWD_REQUIRED
	}
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return errors.New("could not find remember")
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) defaultify(user *User) error {
	if user.Remember != "" {
		return errors.New("could not find remember")
	}

	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err

	}
	if n < 32 {
		return ERR_REMMEMBER_TOO_SHOT
	}
	return nil
}
func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ERR_REMMEMBER_REQUIRED
	}
	return nil
}
func (uv *userValidator) idGreaterThanZero(user *User) error {
	if user.ID <= 0 {
		return ERR_ID_INVALID
	}
	return nil
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)

	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ERR_MAIL_REQUIRED
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ERR_MAIL_IS_N0T_VALID
	}
	return nil
}

func (uv *userValidator) emailsIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)

	switch err {
	case ERR_NOT_FOUND:
		return nil
	case nil:

	default:
		return ERR_MAIL_NOT_EXIST
	}

	if user.ID != existing.ID {
		return ERR_MAIL_IS_TAKEN
	}

	return nil
}

func (uv *userValidator) Create(user *User) error {

	if err := runUserValFuncs(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.defaultify,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailsIsAvail); err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	if err := runUserValFuncs(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.defaultify,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailsIsAvail); err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id string) error {
	if !valid.IsInt(id) {
		return ERR_ID_INVALID
	}
	return uv.UserDB.Delete(id)
}

func (uv *userValidator) ByRemember(token string) (*User, error) {

	if token == "" {
		return nil, ERR_REMMEMBER_REQUIRED
	}

	return uv.UserDB.ByRemember(token)
}

func newUserGorm(db *gorm.DB) (*userGorm, error) {
	return &userGorm{
		db: db,
	}, nil
}

var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id).First(&user)
	err := first(db, &user)
	return &user, err
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User

	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err

}

func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := ug.db.Preload("Rol").Where("remember_hash = ?", rememberHash).First(&user).Error
	return &user, err
}

func (us *userService) InitiateReset(userID string) (string, error) {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return "", ERR_ID_INVALID
	}
	pwr := pwReset{
		UserID: uint(id),
	}
	if err := us.pwResetDB.Create(&pwr); err != nil {
		return "", err
	}
	return pwr.TokenHash, nil
}

func (us *userService) CompleteReset(token, newPw string) (*User, error) {
	pwr, err := us.pwResetDB.ByToken(token)
	if err != nil {
		return nil, err
	}
	if time.Since(pwr.CreatedAt) > (2 * time.Hour) {
		return nil, err
	}
	user, err := us.ByID(pwr.UserID)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(newPw+userPwPPepper))
	if err == nil {
		return nil, ERR_PSSWD_SAME_RESET
	}
	user.Password = newPw
	err = us.Update(user)
	if err != nil {
		return nil, err
	}
	us.pwResetDB.Delete(pwr.ID)

	return user, nil

}

func (us *userService) Authenticate(email, password string) (*User, error) {
	if email == "" {
		return nil, ERR_MAIL_REQUIRED
	}
	if password == "" {
		return nil, ERR_PSSWD_REQUIRED
	}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\]+@[a-z0-9.\-]+\.[a-z]{2,16}$`)
	if !emailRegex.MatchString(email) {
		return nil, ERR_MAIL_IS_N0T_VALID
	}
	email = strings.ToLower(email)
	email = strings.TrimSpace(email)

	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, ERR_MAIL_NOT_EXIST
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ERR_PSSWD_INCORRECT
		default:
			return nil, err
		}
	}

	return foundUser, nil
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return ERR_NOT_FOUND
	default:
		return err
	}
}

func (ug *userGorm) Create(user *User) error {
	err := ug.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (ug *userGorm) Delete(id string) error {
	varInt, err := strconv.Atoi(id)
	if err != nil {
		return ERR_ID_INVALID
	}
	user := User{NeftModel: NeftModel{ID: uint(varInt)}}
	return ug.db.Delete(&user).Error
}

func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

func (ug *userGorm) GetAllUsers() ([]*User, error) {
	var users []*User
	err := ug.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

type User struct {
	NeftModel
	Name         string    `gorm:"not null" json:"username"`
	FullName     string    `json:"full_name"`
	Email        string    `gorm:"not null;unique_index" json:"email"`
	Password     string    `gorm:"-" json:"password"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Remember     string    `gorm:"-" json:"-"`
	RememberHash string    `gorm:"not null;unique_index" json:"-"`
	RolID        int       `gorm:"not null;default:1" json:"rolid"`
	Rol          Rol       `gorm:"foreignkey:RolID" json:"rol"`
	Enabled      int       `gorm:"not null;default:1" json:"activated"`
	Photo        string    `gorm:"default:null" json:"photo"`
	DOB          time.Time `json:"dob"`
}
