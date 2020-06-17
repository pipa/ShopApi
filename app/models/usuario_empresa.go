package models

import "gopkg.in/mgo.v2/bson"

// Represents a movie, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document
type UsuarioEmpresa struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	User        string        `bson:"user" json:"user"`
	EmpresaId   bson.ObjectId `bson:"empresaId" json:"empresaId"`
	Clave		string		  `bson:"clave" json:"clave"`
	Nombre		string		  `bson:"nombre" json:"nombre"`
}
