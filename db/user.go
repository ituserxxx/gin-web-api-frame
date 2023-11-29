package db

type User struct {
	ID     uint `gorm:"primaryKey"`
	Name   string
	Passwd string
	Phone  string
	Url    string
}

func (User) TableName() string {
	return "user"
}
