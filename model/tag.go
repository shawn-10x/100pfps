package model

import (
	"strings"

	"github.com/shawn-10x/100pfps/db"
)

type Tag struct {
	ProfileID uint   `gorm:"primaryKey"`
	Name      string `gorm:"primaryKey"`
}

type TagsAvaliable []struct {
	Name  string
	Count int
}

func GetAvaliableTags() (tags_avaliable TagsAvaliable) {
	db.GetDB().Model(&Tag{}).Select("name, count(profile_id)").Group("name").Order("name").Scan(&tags_avaliable)
	return
}

func StrToTags(tags_str string) []Tag {
	split_tags := strings.Split(tags_str, " ")
	var tags []Tag

	for _, tag := range split_tags {
		tags = append(tags, Tag{
			Name: tag,
		})
	}
	return tags
}
