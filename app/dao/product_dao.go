package dao

import (
	"log"
	"math/rand"
	"time"
	"fmt"

	"ShopAPI/shopapi/models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ProductsDao struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "products"
)

// Establish a connection to database
func (m *ProductsDao) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of products
func (m *ProductsDao) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := db.C(COLLECTION).Find(bson.M{}).All(&products)
	return products, err
}

// Find a products by its id
func (m *ProductsDao) FindById(id string) (models.Product, error) {
	var products models.Product
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&products)
	return products, err
}

// Find a products by its empresaId
func (m *ProductsDao) FindByEmpresaId(id string) ([]models.Product, error) {
	var products []models.Product
	err := db.C(COLLECTION).Find(bson.M{"empresaId": bson.ObjectIdHex(id)}).All(&products)
	return products, err
}

// Find a products by word
func (m *ProductsDao) FindByWord(productsFilter models.Product) ([]models.Product, error) {
	var products []models.Product
	var tags []string

	tags = append(tags, productsFilter.Tags...)
	tags = append(tags, productsFilter.Name)
	err := db.C(COLLECTION).Find(bson.M{"tags": bson.M{"$in": tags}, "name": bson.M{"$ne": productsFilter.Name}}).All(&products)
	return products, err
}

// Find a products by Category
func (m *ProductsDao) FindByCategory(productsFilter string) ([]models.Product, error) {
	var products []models.Product

	err := db.C(COLLECTION).Find(bson.M{"category.name": productsFilter}).All(&products)
	rand.Seed(time.Now().UnixNano())
	shuffle(products)
	// fmt.Println(products)
	return products, err
}

func(m *ProductsDao) FindProductAddCategory(id string, category models.Category) (models.Product, error) {
	var product models.Product
	
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&product)

	
	if err != nil {
		return product, err
	}
	
	product.Categories = append(product.Categories, category)
	
	e := db.C(COLLECTION).UpdateId(product.ID, &product)

	if e != nil {
		fmt.Println(e)
		return product, e
	}
	
	erro := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&product)
	return product, erro
}


func(m *ProductsDao) UpdateProduct(prod models.Product) (models.Product, error) {
	
	
	e := db.C(COLLECTION).UpdateId(prod.ID,&prod)

	if e != nil {
		fmt.Println(e)
		return prod, e
	}
	
	erro := db.C(COLLECTION).FindId(prod.ID).One(&prod)
	return prod, erro
}


func shuffle(r []models.Product) {
	for i := len(r) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		r[i], r[j] = r[j], r[i]
	}
}

// Insert a products into database
func (m *ProductsDao) Insert(products models.Product) error {
	i := bson.NewObjectId()

	products.ID = i
	fmt.Println(products)
	err := db.C(COLLECTION).Insert(&products)
	return err
}

// Delete an existing products
func (m *ProductsDao) Delete(products models.Product) error {
	err := db.C(COLLECTION).Remove(&products)
	return err
}

// Update an existing products
func (m *ProductsDao) Update(products models.Product) error {
	err := db.C(COLLECTION).UpdateId(products.ID, &products)
	return err
}
