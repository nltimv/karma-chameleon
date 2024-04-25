package db

type User struct {
	ID     uint `gorm:"primaryKey"`
	UserId string
	TeamId string
	Karma  int
}

type Group struct {
	ID      uint `gorm:"primaryKey"`
	GroupId string
	TeamId  string
	Karma   int
}
