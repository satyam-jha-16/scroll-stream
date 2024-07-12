package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `gorm:"size:255;not null"`
	Email        string `gorm:"size:255;not null;unique"`
	Password     string `gorm:"size:255;not null"`
	ProfileImage string `gorm:"size:255"`
	LikedVideos  []Video `gorm:"many2many:user_liked_videos;"`
}

type Video struct {
	gorm.Model
	Title        string `gorm:"size:255;not null"`
	Description  string `gorm:"type:text"`
	VideoURL     string `gorm:"size:255;not null"`
	ImageURL     string `gorm:"size:255"`
	Likes        int    `gorm:"default:0"`
	PublishedBy  uint   
	Publisher    User   `gorm:"foreignKey:PublishedBy"`
	LikedBy      []User `gorm:"many2many:user_liked_videos;"`
}
