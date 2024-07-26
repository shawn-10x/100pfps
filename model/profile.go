package model

import (
	"fmt"
	"net"
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

func GetProfile(id uint) (profile *Profile) {
	profile = new(Profile)
	err := db.GetDB().Preload("Tags").Where("id = ?", id).Take(profile).Error
	if err != nil {
		profile = nil
	}
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
	return nil
}

func ExistsProfileWithIP(ip net.IP) (exists bool) {
	var count int64
	db.GetDB().Model(&Profile{}).Where("ip = ?", ip.String()).Count(&count)
	return count > 0
}

func (profile *Profile) Insert() (err error) {
	return db.GetDB().Preload("Tags").Create(profile).Error
}

func (profile *Profile) GetIPNet() net.IPNet {
	netip := *profile.Ip.IPNet
	if netip.IP.To4() != nil {
		netip.Mask = createMask(32)
	} else {
		netip.Mask = createMask(70)
	}
	return netip

}

func (u *Profile) BeforeCreate(db *gorm.DB) (err error) {
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
