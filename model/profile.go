package model

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jackc/pgtype"
	"github.com/shawn-10x/100pfps/db"
)

type Profile struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
	Tags        []Tag
	Ip          pgtype.Inet `gorm:"uniqueIndex;type:inet;not null"`
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

func GetProfileImg(profile *Profile) string {
	return fmt.Sprintf("imgs/%d.png", profile.ID)
}

func DeleteProfile(profile *Profile) (err error) {
	db := db.GetDB()
	if err = db.Delete(profile).Error; err != nil {
		return
	}
	if err = os.Remove(GetProfileImg(profile)); err != nil {
		return fmt.Errorf("Error deleting image")
	}
	return nil
}

func ExistsProfileWithIP(ip net.IP) (exists bool) {
	var count int64
	db.GetDB().Model(&Profile{}).Where("ip = ?", ip.String()).Count(&count)
	return count > 0
}

func InsertProfile(profile *Profile) (err error) {
	db := db.GetDB()

	var count int64
	db.Model(&Profile{}).Count(&count)
	if count >= 99 {
		var profile Profile
		db.Order("id ASC").First(&profile)
		if err = DeleteProfile(&profile); err != nil {
			return
		}
	}
	return db.Preload("Tags").Create(profile).Error
}
