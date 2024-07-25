package model

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jackc/pgtype"
	"github.com/shawn-10x/100pfps/db"
	"gorm.io/gorm"
)

type Profile struct {
	ID          uint        `gorm:"primaryKey"`
	Name        string      `gorm:"not null"`
	Description string      `gorm:"not null"`
	Tags        []Tag       `gorm:"constraint:OnDelete:CASCADE"`
	Ip          pgtype.Inet `gorm:"uniqueIndex;type:inet;not null"`
	Thumbnail   []byte      `gorm:"not null"`
	Image       []byte      `gorm:"not null"`
	CreatedAt   time.Time
}

func GetProfiles(tag *string) (profiles []Profile) {
	tx := db.GetDB()
	if tag == nil || *tag == "any" {
		tx = tx.Preload("Tags").Find(&profiles)
	} else {
		tx = tx.Joins("JOIN tags ON tags.profile_id = profiles.id AND tags.name = ?", *tag)
	}
	tx.Order("id DESC").Find(&profiles)
	return
}

func (profile *Profile) GetProfileImg() string {
	return fmt.Sprintf("imgs/%d.png", profile.ID)
}

func (profile *Profile) Delete() (err error) {
	db := db.GetDB()
	if err = db.Delete(profile).Error; err != nil {
		return
	}
	if err = os.Remove(profile.GetProfileImg()); err != nil {
		return fmt.Errorf("Error deleting image")
	}
	return nil
}

func ExistsProfileWithIP(ip net.IP) (exists bool) {
	var count int64
	db.GetDB().Model(&Profile{}).Where("ip = ?", ip.String()).Count(&count)
	return count > 0
}

func (profile *Profile) Insert() (err error) {
	db := db.GetDB()
	return db.Preload("Tags").Create(profile).Error
}

func (u *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	db := db.GetDB()

	var count int64
	db.Model(&Profile{}).Count(&count)
	if count >= 99 {
		var profile Profile
		db.Order("id ASC").First(&profile)
		if err = profile.Delete(); err != nil {
			return
		}
	}
	return
}
