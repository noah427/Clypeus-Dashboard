package main


import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type GuildSettings struct {
	gorm.Model
	ID         string `json:"ID"`
	AntinukeON bool   `json:"AntinukeON"`
	AntispamON bool   `json:"AntispamON"`
	MuteRole   string
	TrustedID  []User `gorm:"foreignKey:ID"`
	Prefix     string
}

type User struct {
	gorm.Model
	ID       string
	Username string
}

func loadDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&GuildSettings{})
}
