package models

type Course struct {
	ID   uint   `gorm:"primaryKey" json:"-"`
	Guid string `gorm:"size:55;not null"`
	name string `gorm:"primaryKey" json:"-"`
}
