package models

type Mail struct {
	To               Email  `json:"to"`
	Subject          string `json:"subject"`
	PlainTextContent string `json:"plainTextContent"`
	HtmlContent      string `json:"htmlContent"`
}

type Email struct {
	Name     string `json:"name"`
	UserMail string `json:"userMail"`
}
