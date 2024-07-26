package model

import (
	"crypto/sha512"

	"github.com/shawn-10x/100pfps/db"
	"github.com/shawn-10x/100pfps/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminRole uint8

const (
	Owner AdminRole = iota
	Moderator
	Helper
)

type Admin struct {
	User     string `gorm:"primaryKey"`
	Role     AdminRole
	Password string
	Token    string `gorm:"uniqueIndex"` // No jwt please
}

func (admin *Admin) GetRoleStr() string {
	switch admin.Role {
	case Owner:
		return "Owner"
	case Moderator:
		return "Moderator"
	case Helper:
		return "Helper"
	}
	return "???"
}

func (admin *Admin) HashPassword() (err error) {
	hash := sha512.Sum384([]byte(admin.Password))
	pw, err := bcrypt.GenerateFromPassword(hash[:], bcrypt.DefaultCost)
	admin.Password = string(pw)
	return
}

func (admin *Admin) CheckPassword(pw string) bool {
	hash := sha512.Sum384([]byte(pw))
	return bcrypt.CompareHashAndPassword([]byte(admin.Password), hash[:]) == nil
}

func (admin *Admin) NewToken() {
	admin.Token = utils.GenerateRandomString(32)
}

func (admin *Admin) Save() (err error) {
	return db.GetDB().Save(admin).Error
}

func (admin *Admin) Create() (err error) {
	return db.GetDB().Create(admin).Error
}

func (admin *Admin) CreateIfNotExists() (err error) {
	if GetAdmin(admin.User) != nil {
		return
	}
	return admin.Create()
}

func GetAdmin(user string) (admin *Admin) {
	admin = new(Admin)
	err := db.GetDB().Where("admins.user = ?", user).Take(admin).Error
	if err != nil {
		admin = nil
	}
	return
}

func IsAdmin(token string) bool {
	var count int64
	db.GetDB().Model(&Admin{}).Where("token = ?", token).Count(&count)
	return count > 0
}

func GetAdminByToken(token string) (admin *Admin) {
	admin = new(Admin)
	err := db.GetDB().Where("token = ?", token).Take(admin).Error
	if err != nil {
		admin = nil
	}
	return
}

func GetAdmins(user string) (admins []Admin) {
	db.GetDB().Find(&admins)
	return
}

func (admin *Admin) BeforeCreate(db *gorm.DB) (err error) {
	admin.HashPassword()
	admin.NewToken()
	return
}
