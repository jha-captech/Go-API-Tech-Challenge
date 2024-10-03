package models

type PersonType string

const (
	Professor PersonType = "professor"
	Student   PersonType = "student"
)

type Person struct {
	ID        uint   `gorm:"primaryKey" json:"-"`
	Guid      string `gorm:"size:55;not null"`
	FirstName string `gorm:"size:255;not null"`
	LastName  string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255;not null"`
	Age       uint
	Type      string   `gorm:"size:255;not null"`
	Courses   []Course `gorm:"many2many:person_course;"`
}

func (Person) TableName() string {
	return "person"
}

type PersonInput struct {
	FirstName   string
	LastName    string
	Email       string
	Age         uint
	Type        string
	CourseGuids []string
}
