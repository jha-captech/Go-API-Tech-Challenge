package models

type PersonCourse struct {
	PersonID uint `gorm:"primaryKey" json:"-"`
	CourseID uint `gorm:"primaryKey" json:"-"`
}
