package dao

import (
	"log"	
	"ShopAPI/shopapi/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UsuarioDao struct {
	Server   string
	Database string
}

var usuariodb *mgo.Database

const (
	UsuarioCOLLECTION = "usuarios_empresa"
)


// User ...
// Custom object which can be stored in the claims
type User struct {
	Username string `json:"username"`
	EmpresaId string `json:"empresaId"`
	Nombre	 string  `json:"nombre"`
}

// AuthToken ...
// This is what is retured to the user
type AuthToken struct {
	TokenType string `json:"token_type"`
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

// AuthTokenClaim ...
// This is the cliam object which gets parsed from the authorization header
type AuthTokenClaim struct {
	*jwt.StandardClaims
	User
}

// Establish a connection to database
func (m *UsuarioDao) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	usuariodb = session.DB(m.Database)
}

// Find list of products
func (m *UsuarioDao) FindAll() ([]models.UsuarioEmpresa, error) {
	var usuarios []models.UsuarioEmpresa
	err := usuariodb.C(UsuarioCOLLECTION).Find(bson.M{}).All(&usuarios)
	return usuarios, err
}

// Find a products by its id
func (m *UsuarioDao) FindById(id string) (models.UsuarioEmpresa, error) {
	var usuarios models.UsuarioEmpresa
	err := usuariodb.C(UsuarioCOLLECTION).FindId(bson.ObjectIdHex(id)).One(&usuarios)
	return usuarios, err
}

// func empresaShuffle(r []models.UsuarioEmpresa) {
// 	for i := len(r) - 1; i > 0; i-- {
// 		j := rand.Intn(i + 1)
// 		r[i], r[j] = r[j], r[i]
// 	}
// }

// Insert a products into database
func (m *UsuarioDao) Insert(usuario models.UsuarioEmpresa) error {
	err := usuariodb.C(UsuarioCOLLECTION).Insert(&usuario)
	return err
}

//Validating
func (m *UsuarioDao) Validate(usuario models.UsuarioEmpresa) ( string , error) {
	query := bson.M{"user":  bson.RegEx{ usuario.User , ""}}
	var usuarioDB models.UsuarioEmpresa
	var response string
	err:= usuariodb.C(UsuarioCOLLECTION).Find(query).One(&usuarioDB)

	// log.Println(usuarioDB)
	response = ""
	if usuario.User == usuarioDB.User && usuario.Clave == usuarioDB.Clave && usuario.EmpresaId == usuarioDB.EmpresaId {
		//response = true
		expiresAt := time.Now().Add(time.Hour * 7).Unix()

		token := jwt.New(jwt.SigningMethodHS256)

		token.Claims = &AuthTokenClaim{
			&jwt.StandardClaims{
				ExpiresAt: expiresAt,
			},
			User{usuario.User, usuario.EmpresaId.Hex(), usuarioDB.Nombre},
		}

		tokenString, error := token.SignedString([]byte("Api2019"))
		if error != nil {
			response = ""
		}

		response = tokenString
	}

	return response, err
}

// Delete an existing products
func (m *UsuarioDao) Delete(usuario models.UsuarioEmpresa) error {
	err := usuariodb.C(UsuarioCOLLECTION).Remove(&usuario)
	return err
}

// Update an existing products
func (m *UsuarioDao) Update(usuario models.UsuarioEmpresa) error {
	err := usuariodb.C(UsuarioCOLLECTION).UpdateId(usuario.ID, &usuario)
	return err
}
