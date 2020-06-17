package dao

import (
	"log"
	"math/rand"
	
	"ShopAPI/shopapi/models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type EmpresaDao struct {
	Server   string
	Database string
}

var empresadb *mgo.Database

const (
	EmpresaCOLLECTION = "empresas"
)

// Establish a connection to database
func (m *EmpresaDao) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	empresadb = session.DB(m.Database)
}

// Find list of products
func (m *EmpresaDao) FindAll() ([]models.Empresa, error) {
	var empresas []models.Empresa
	err := empresadb.C(EmpresaCOLLECTION).Find(bson.M{}).All(&empresas)
	return empresas, err
}

// Find a products by its id
func (m *EmpresaDao) FindById(id string) (models.Empresa, error) {
	var empresas models.Empresa
	err := empresadb.C(EmpresaCOLLECTION).FindId(bson.ObjectIdHex(id)).One(&empresas)
	return empresas, err
}

// Find a products by its id
func (m *EmpresaDao) FindByCompany(company string) (models.Empresa, error) {
	var empresas models.Empresa
	query := bson.M{"name":  bson.RegEx{ company, ""}}

	err := empresadb.C(EmpresaCOLLECTION).Find(query).One(&empresas)
	return empresas, err
}

func empresaShuffle(r []models.Empresa) {
	for i := len(r) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		r[i], r[j] = r[j], r[i]
	}
}

// Insert a products into database
func (m *EmpresaDao) Insert(empresas models.Empresa) error {
	err := empresadb.C(EmpresaCOLLECTION).Insert(&empresas)
	return err
}

// Delete an existing products
func (m *EmpresaDao) Delete(empresas models.Empresa) error {
	err := empresadb.C(EmpresaCOLLECTION).Remove(&empresas)
	return err
}

// Update an existing products
func (m *EmpresaDao) Update(empresas models.Empresa) error {
	err := empresadb.C(EmpresaCOLLECTION).UpdateId(empresas.ID, &empresas)
	return err
}
