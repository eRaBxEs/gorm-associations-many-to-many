package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CREATE TABLES WITH THE FOLLOWING DETAILS: Movies, Actors, Filmography

/* Sample Data:
   <Movie>: 					<List of Actors>
   "Iron Man": 					Robert Downey Jr.
   "Avengers": 					Robert Downey Jr., Chris Evans, Scarlett Johansson
   "Black Panther": 			Chadwick Boseman
   "Avengers Infinity War": 	Robert Downey Jr., Chris Evans, Scarlett Johansson, and Chadwick Boseman
   "Sherlock Holmes": 			Robert Downey Jr.
   "Lost in Translation": 		Scarlett Johansson
   "Marriage Story": 			Scarlett Johansson
*/

type Movie struct {
	gorm.Model
	Name   string
	Actors []Actor `gorm:"many2many:filmography;"`
}

type Actor struct {
	gorm.Model
	Name   string
	Movies []Movie `gorm:"many2many:filmography;"` // added in order to enable one be able to specify the movies an actor has acted in
}

var DB *gorm.DB

func connectDatabase() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	dsn := "henry_dev:devdba_user@tcp(127.0.0.1:3307)/gorm_testdb?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic("Failed to connect to databse!")
	}

	DB = database
}

func dbMigrate() {
	DB.AutoMigrate(&Movie{}, &Actor{})
}

func main() {
	connectDatabase()
	dbMigrate()

	var movie Movie
	DB.Where("name = ?", "Avengers Infinity War").Preload("Actors").First(&movie)
	fmt.Printf("Movie: %s\n\n", movie.Name)

	fmt.Println("Actors:")
	for _, actor := range movie.Actors {
		fmt.Printf("%s\n", actor.Name)
	}

	// similarly if we are to get all the movies that Robert Downey Jr. has acted in
	var actor Actor
	DB.Where("name = ?", "Robert Downey Jr.").Preload("Movies").First(&actor)
	fmt.Println("Actor: " + actor.Name)

	fmt.Printf("Movies:")

	for _, element := range actor.Movies {
		fmt.Printf("%v\n", element.Name)
	}

}
