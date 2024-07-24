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

func GetProfiles(tag *string) (profiles []Profile, err error) {
	tx := db.GetDB()
	if tag == nil || *tag == "any" {
		tx = tx.Preload("Tags").Find(&profiles)
	} else {
		tx = tx.Joins("JOIN tags ON tags.profile_id = profiles.id AND tags.name = ?", *tag)
	}
	err = tx.Order("id DESC").Find(&profiles).Error
	return
}

func GetProfileImg(profile *Profile) string {
	return fmt.Sprintf("/imgs/%d.png", profile.ID)
}

func DeleteProfile(profile *Profile) (err error) {
	db := db.GetDB()
	if err = db.Delete(profile).Error; err != nil {
		return
	}
	if err = os.Remove(GetProfileImg(profile)); err != nil {
		return
	}
	return nil
}

func ExistsProfileWithIP(ip net.IP) (exists bool, err error) {
	var count int64
	err = db.GetDB().Model(&Profile{}).Where("ip = ?", ip.String()).Count(&count).Error
	return count > 0, err
}

func InsertProfile(profile *Profile) (err error) {
	db := db.GetDB()

	var count int64
	db.Model(&Profile{}).Count(&count)
	if count >= 99 {
		var profile Profile
		if err = db.Order("id ASC").First(&profile).Error; err != nil {
			return
		}
		if err = DeleteProfile(&profile); err != nil {
			return
		}
	}
	return db.Preload("Tags").Create(profile).Error
}
