package dao

import (
	"log"

	"ShopAPI/shopapi/models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CategoryDao struct {
	Server   string
	Database string
}

var categoryDb *mgo.Database

const (
	categoryCOLLECTION = "categories"
)

// Establish a connection to database
func (m *CategoryDao) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	categoryDb = session.DB(m.Database)
}

// Find list of categories
func (m *CategoryDao) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := categoryDb.C(categoryCOLLECTION).Find(bson.M{}).All(&categories)
	return categories, err
}


func (m *CategoryDao) FindAllByCompany(id string) ([]models.Category, error) {
	var categories []models.Category
	err := categoryDb.C(categoryCOLLECTION).Find(bson.M{"empresaId": bson.ObjectIdHex(id)}).All(&categories)
	return categories, err
}


// Find a categories by its id
func (m *CategoryDao) FindById(id string) (models.Category, error) {
	var categories models.Category
	err := categoryDb.C(categoryCOLLECTION).FindId(bson.ObjectIdHex(id)).One(&categories)
	return categories, err
}

// Insert a categories into database
func (m *CategoryDao) Insert(categories models.Category) error {
	err := categoryDb.C(categoryCOLLECTION).Insert(&categories)
	return err
}

// Delete an existing categories
func (m *CategoryDao) Delete(categories models.Category) error {
	err := categoryDb.C(categoryCOLLECTION).Remove(&categories)
	return err
}

// Update an existing categories
func (m *CategoryDao) Update(categories models.Category) error {
	err := categoryDb.C(categoryCOLLECTION).UpdateId(categories.ID, &categories)
	return err
}
