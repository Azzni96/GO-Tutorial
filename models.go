package main

import "time"

type Article struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string    `json:"title"  gorm:"type:varchar(200);not null"`
	Desc      string    `json:"desc"   gorm:"type:varchar(500)"`
	Content   string    `json:"content" gorm:"type:text"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
