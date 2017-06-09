package parsers

import "regexp"

type EmailData struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	EmailAddress string `json:"email_address"`
	Phone        string `json:"phone"`
	Status       string `json:"status"`
	EnquiryType  string `json:"enquiry"`
	ImageURL     string `json:"image_url"`
	Address      string `json:"address"`
	Description  string `json:"description"`
	Price        string `json:"price"`
	Link         string `json:"link"`
	CreatedAt    string `json:"created_at"`
	HTML         string `json:"html"`
}

func (ed *EmailData) IsFailed() bool {
	return (ed.Name == "" || ed.EmailAddress == "" || ed.Address == "")
}

func Parse(from string, body string, emailData *EmailData) {
	//switch from {
	//case "lavavrik@yandex.ru":
	ParseZoopla(body, emailData)
	//}
}

func RemoveHTML(input string) string {
	re := regexp.MustCompile(`<.*?>`)
	return string(re.ReplaceAll([]byte(input), []byte(""))[:])
}
