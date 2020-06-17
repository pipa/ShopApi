package models

import "gopkg.in/mgo.v2/bson"

// Represents a movie, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document
type Product struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Handle		string		  `bson:"handle" json:"handle"`
	Description string        `bson:"description" json:"description"`
	Categories  []Category 	  `bson:"categories" json:"categories"`
	Tags        []string      `bson:"tags" json:"tags"`
	FeaturedImageId int		  `bson:"featuredImageId" json:"featuredImageId"`
	Images		[]Image		  `bson:"images" json:"images"`
	PriceTaxExcl float64        `bson:"priceTaxExcl" json:"priceTaxExcl"`
	PriceTaxIncl float64		  `bson:"priceTaxIncl" json:"priceTaxIncl"`	
	TaxRate		float64		  `bson:"taxRate" json:"taxRate"`
	ComparedPrice float64		  `bson:"comparedPrice" json:"comparedPrice"`
	Quantity     int64        `bson:"quantity" json:"quantity"`
	Sku			string		  `bson:"sku" json:"sku"`
	Width		string		  `bson:"width" json:"width"`
	Height		string		  `bson:"height" json:"height"`
	Depth		string		  `bson:"depth" json:"depth"`
	Weight		string		  `bson:"weight" json:"weight"`	
	ExtraShippingFee	float64	`bson:"extraShippingFee" json:"extraShippingFee"`	
	Active		bool		  `bson:"active" json:"active"`	
	
	EmpresaId	bson.ObjectId `bson:"empresaId" json:"empresaId"`
}


type Image struct {
	ID		int64	`bson:"_id" json:"id"`
	Url		string			`bson:"url" json:"url"`
	Type	string			`bson:"type" json:"type"`
}