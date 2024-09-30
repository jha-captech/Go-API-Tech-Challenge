package models

type PersonCourse struct {
	PersonID uint `gorm:"primaryKey"`
	CourseID uint `gorm:"primaryKey"`
}

func (PersonCourse) TableName() string {
	return "person_course"
}
