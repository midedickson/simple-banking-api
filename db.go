package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func connectToDB() {
	d, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Println("Connected to database sucessfully")
	DB = d
}

func autoMigrate() {
	log.Println("Auto Migrating Models...")
	err := DB.AutoMigrate(&Transaction{})
	if err != nil {
		panic(err)
	}
	log.Println("Migrated DB Successfully")
}
