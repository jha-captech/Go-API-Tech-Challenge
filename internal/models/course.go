package models

type CourseInput struct {
	Name string
}

type Course struct {
	ID   uint   `gorm:"primaryKey" json:"-"`
	Guid string `gorm:"size:55;not null"`
	Name string `gorm:"size:255;not null"`
}

func (Course) TableName() string {
	return "course"
}
