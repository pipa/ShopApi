package models

import "gopkg.in/mgo.v2/bson"

type Category struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
	EmpresaId bson.ObjectId `bson:"empresaId" json:"empresaId"`
}
