package model

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"os"
)

type Author struct {
	gorm.Model
	AuthorID   int    `json:"ID"`
	AuthorName string `json:"Name"`
}

func (a *Author) ToString() string {
	return fmt.Sprintf("AuthorID: %d, AuthorName: %s", a.AuthorID, a.AuthorName)
}

func GetAllAuthorsFromJson() []Author {
	authors := []Author{}

	jsonFile, err := os.Open("authors.json")

	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal("Patates while opening json: ", err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	//Parse to Authors struct
	err = json.Unmarshal(byteValue, &authors)

	//fmt.Println(byteValue)
	if err != nil {
		log.Fatal("Patates while unmarshal json: ", err)
	}

	return authors
}

func AuthorDBMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Author{})

	if err != nil {
		fmt.Println("Migration Error")
		return err
	}

	return nil
}
