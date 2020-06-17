package main

import (
	// "crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pipa/ShopAPI/config"
	"github.com/pipa/ShopAPI/dao"
	"github.com/pipa/ShopAPI/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/gomail.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
)

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender     = "cotizaciones@xn--elgreas-8za.com"
	SenderName = "Greñas Shop"

	// Replace recipient@example.com with a "To" address. If your account
	// is still in the sandbox, this address must be verified.
	Recipient = "josedavidgale@gmail.com"

	// Replace SmtpUser with your Amazon SES SMTP user name.
	SmtpUser = "AKIAIWIJDTYWYRQ65C7A"

	// Replace SmtpPass with your Amazon SES SMTP password.
	SmtpPass = "Ahqpsm/QcSEHJIKgM/8IHoOzbewaPyBtZVik2dOMT1P2"

	// The name of the configuration set to use for this message.
	// If you comment out or remove this variable, you will also need to
	// comment out or remove the header below.
	// ConfigSet = "ConfigSet"

	// If you're using Amazon SES in an AWS Region other than US West (Oregon),
	// replace email-smtp.us-west-2.amazonaws.com with the Amazon SES SMTP
	// endpoint in the appropriate region.
	Host = "email-smtp.us-east-1.amazonaws.com"
	Port = 587

	// The subject line for the email.
	Subject = "Amazon SES Test (Gomail)"

	// The HTML body for the email.
	HtmlBody = "<html><head><title>SES Sample Email</title></head><body>" +
		"<h1>Amazon SES Test Email (Gomail)</h1>" +
		"<p>This email was sent with " +
		"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using " +
		"the <a href='https://github.com/go-gomail/gomail/'>Gomail " +
		"package</a> for <a href='https://golang.org/'>Go</a>.</p>" +
		"</body></html>"

	//The email body for recipients with non-HTML email clients.
	TextBody = "This email was sent with Amazon SES using the Gomail package."

	// The tags to apply to this message. Separate multiple key-value pairs
	// with commas.
	// If you comment out or remove this variable, you will also need to
	// comment out or remove the header on line 80.
	Tags = "genre=test,genre2=test2"

	// The character encoding for the email.
	CharSet = "UTF-8"
)

var Config = config.Config{}
var DAO = dao.ProductsDao{}
var CategortDAO = dao.CategoryDao{}
var EmpresaDAO = dao.EmpresaDao{}
var UsuarioDAO = dao.UsuarioDao{}



// User ...
// Custom object which can be stored in the claims
type User struct {
	Username string `json:"username"`
	EmpresaId string `json:"empresaId"`
	Nombre	 string `json:"nombre"`
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

// ErrorMsg ...
// Custom error object
type ErrorMsg struct {
	Message string `json:"message"`
}

// Create the JWT key used to create the signature
var jwtKey = []byte("Api2019")

func AllProductsEndPoint(w http.ResponseWriter, r *http.Request) {
	products, err := DAO.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, products)
}

func FindProductEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	product, err := DAO.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	respondWithJson(w, http.StatusOK, product)
}

func FindProductByCompanyEndpoint(w http.ResponseWriter, r *http.Request) {
	
	decoded := context.Get(r, "decoded")
    var user User
    mapstructure.Decode(decoded.(jwt.MapClaims), &user)
    // json.NewEncoder(w).Encode(user)	

	product, err := DAO.FindByEmpresaId(user.EmpresaId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	respondWithJson(w, http.StatusOK, product)
}

func FindProductWordEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var productFilter models.Product
	if err := json.NewDecoder(r.Body).Decode(&productFilter); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	product, err := DAO.FindByWord(productFilter)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	respondWithJson(w, http.StatusOK, product)
}

func FindProductCategoryEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// if err := json.NewDecoder(r.Body).Decode(&productFilter); err != nil {
	// 	respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	// 	return
	// }
	params := mux.Vars(r)
	product, err := DAO.FindByCategory(params["category"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	respondWithJson(w, http.StatusOK, product)
}

func AllCompanyEndPoint(w http.ResponseWriter, r *http.Request) {
	empresas, err := EmpresaDAO.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, empresas)
}

func FindCompaniesEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	empresas, err := EmpresaDAO.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	respondWithJson(w, http.StatusOK, empresas)
}

func FindCompanyEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	empresas, err := EmpresaDAO.FindByCompany(params["company"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, empresas)
}

func CreateProductEndPoint(w http.ResponseWriter, r *http.Request) {
	decoded := context.Get(r, "decoded")
    var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)

	
	defer r.Body.Close()
	var product models.Product
	// fmt.Println(r.Body)

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		fmt.Println(err);
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	product.ID = bson.NewObjectId()
	product.EmpresaId = bson.ObjectIdHex(user.EmpresaId)
	if err := DAO.Insert(product); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, product)
}

func CreateCategoryEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	category.ID = bson.NewObjectId()
	// log.Println(category)
	if err := CategortDAO.Insert(category); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, category)
}

func CreateEmpresaEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var empresa models.Empresa
	if err := json.NewDecoder(r.Body).Decode(&empresa); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	empresa.ID = bson.NewObjectId()
	// log.Println(category)
	if err := EmpresaDAO.Insert(empresa); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, empresa)
}

func CreateUsuarioEmpresaEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var usuario models.UsuarioEmpresa
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	usuario.ID = bson.NewObjectId()
	// log.Println(category)
	if err := UsuarioDAO.Insert(usuario); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, usuario)
}

func ValidateUsuarioEmpresaEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var usuario models.UsuarioEmpresa
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	response, err := UsuarioDAO.Validate(usuario)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, response)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func sendMail(w http.ResponseWriter, r *http.Request) {
	var email models.Mail

	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// //log.Println(email)
	// from := mail.NewEmail("Jose Gale", "josedavidgale@gmail.com")
	// subject := email.Subject
	// to := mail.NewEmail(email.To.Name, email.To.UserMail)
	// plainTextContent := email.PlainTextContent
	// htmlContent := email.HtmlContent
	// message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	// client := sendgrid.NewSendClient("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")
	// response, err := client.Send(message)
	// if err != nil {
	// 	log.Println(err)
	// } else {

	// 	respondWithJson(w, response.StatusCode, "Email ")
	// 	// fmt.Println(response.StatusCode)
	// 	// fmt.Println(response.Body)
	// 	// fmt.Println(response.Headers)
	// }

	// Create a new message.
	m := gomail.NewMessage()

	// Set the main email part to use HTML.
	m.SetBody("text/html", email.HtmlContent)

	// Set the alternative part to plain text.
	// m.AddAlternative("text/plain", TextBody)

	// Construct the message headers, including a Configuration Set and a Tag.
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(Sender, SenderName)},
		"To":      {email.To.UserMail},
		"Subject": {email.Subject},
		// Comment or remove the next line if you are not using a configuration set
		// "X-SES-CONFIGURATION-SET": {ConfigSet},
		// Comment or remove the next line if you are not using custom tags
		"X-SES-MESSAGE-TAGS": {Tags},
	})

	// Send the email.
	d := gomail.NewPlainDialer(Host, Port, SmtpUser, SmtpPass)

	// Display an error message if something goes wrong; otherwise,
	// display a message confirming that the message was sent.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email sent!")
		respondWithJson(w, 200, "Email ")
	}

}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	Config.Read()

	DAO.Server = Config.Server
	DAO.Database = Config.Database
	DAO.Connect()

	CategortDAO.Server = Config.Server
	CategortDAO.Database = Config.Database
	CategortDAO.Connect()

	EmpresaDAO.Server = Config.Server
	EmpresaDAO.Database = Config.Database
	EmpresaDAO.Connect()

	UsuarioDAO.Server = Config.Server
	UsuarioDAO.Database = Config.Database
	UsuarioDAO.Connect()

	
}


func validateTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// fmt.Println("Aqui")
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		// fmt.Println(authorizationHeader)
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte(jwtKey), nil
				})
				if error != nil {
					json.NewEncoder(w).Encode(ErrorMsg{Message: error.Error()})
					return
				}
				if token.Valid {

					var user User
					mapstructure.Decode(token.Claims, &user)
					
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					json.NewEncoder(w).Encode(ErrorMsg{Message: "Invalid authorization token"})
				}
			} else {
				json.NewEncoder(w).Encode(ErrorMsg{Message: "Invalid authorization token"})
			}
		} else {
			json.NewEncoder(w).Encode(ErrorMsg{Message: "An authorization header is required"})
		}
	})
}

func sendClaims(w http.ResponseWriter, r *http.Request) {
	
	decoded := context.Get(r, "decoded")
    var user User
    mapstructure.Decode(decoded.(jwt.MapClaims), &user)
    // json.NewEncoder(w).Encode(user)	

	respondWithJson(w, http.StatusOK, user)
}


func AllCategoriesCompany(w http.ResponseWriter, r *http.Request) {
	
	decoded := context.Get(r, "decoded")
    var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)


	categories, err := CategortDAO.FindAllByCompany(user.EmpresaId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Empresa ID")
		return
	}
	respondWithJson(w, http.StatusOK, categories)
}

func FindProductAddCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	decoded := context.Get(r, "decoded")
    var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)

	// fmt.Println(params["id"])

	defer r.Body.Close()
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	products, err := DAO.FindProductAddCategory(params["id"], category)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Empresa ID")
		return
	}
	respondWithJson(w, http.StatusOK, products)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) { 
	decoded := context.Get(r, "decoded")
    var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)


	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	product.EmpresaId = bson.ObjectIdHex(user.EmpresaId)

	products, err := DAO.UpdateProduct(product)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Empresa ID")
		return
	}
	respondWithJson(w, http.StatusOK, products)
}



func main() {

	allowedHeaders := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*", "http://localhost:4200/"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	r := mux.NewRouter()
	r.HandleFunc("/products", AllProductsEndPoint).Methods("GET")
	r.HandleFunc("/products", validateTokenMiddleware(CreateProductEndPoint)).Methods("POST")
	r.HandleFunc("/products/{id}", FindProductEndpoint).Methods("GET")
	r.HandleFunc("/products/{id}/AddCategory", validateTokenMiddleware(FindProductAddCategory)).Methods("POST")
	r.HandleFunc("/products/where", FindProductWordEndpoint).Methods("POST")
	r.HandleFunc("/products/category/", FindProductCategoryEndpoint ).Methods("GET")
	r.HandleFunc("/products/company/", validateTokenMiddleware(FindProductByCompanyEndpoint)).Methods("GET")
	r.HandleFunc("/sendMail", sendMail).Methods("POST")
	r.HandleFunc("/categories", CreateCategoryEndPoint).Methods("POST")
	r.HandleFunc("/categories", validateTokenMiddleware(AllCategoriesCompany)).Methods("GET")
	r.HandleFunc("/companies", AllCompanyEndPoint).Methods("GET")
	r.HandleFunc("/company/{id}", FindCompaniesEndpoint).Methods("GET")
	r.HandleFunc("/company/byName/{company}", FindCompanyEndpoint).Methods("GET")
	r.HandleFunc("/company", CreateEmpresaEndPoint).Methods("POST")
	r.HandleFunc("/usuario", CreateUsuarioEmpresaEndPoint).Methods("POST")
	r.HandleFunc("/usuario/validate", ValidateUsuarioEmpresaEndPoint).Methods("POST")
	r.HandleFunc("/usuario/token", validateTokenMiddleware(sendClaims)).Methods("GET")
	r.HandleFunc("/products/update", validateTokenMiddleware(UpdateProduct)).Methods("POST")


	// server := &http.Server{
	// 	Addr:      ":4000",
	// 	TLSConfig: &tls.Config{},
	// 	Handler:   handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(r),
	// }

	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(allowedHeaders,  allowedOrigins, allowedMethods)(r)))

	// Start Server with SSl log.Fatal(server.ListenAndServeTLS("/etc/letsencrypt/live/shopapi.xn--elgreas-8za.com/fullchain.pem", "/etc/letsencrypt/live/shopapi.xn--elgreas-8za.com/privkey.pem"))
}
