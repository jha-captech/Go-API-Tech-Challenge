package models

type Course struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

func (Course) TableName() string {
	return "course"
}
