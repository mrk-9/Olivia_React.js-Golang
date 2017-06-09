package parsers

import (
	"regexp"
)

func ParseZoopla(body string, emailData *EmailData) {
	userInfoRe := regexp.MustCompile(`<td .*?>(Name|Telephone|Email|Personal status|Type of enquiry):[\s\S]+?<img .*?>[\s\S]+?<td .*?>(.+?)</td>`)

	for _, match := range userInfoRe.FindAllStringSubmatch(body, -1) {
		key, value := match[1], match[2]

		switch key {
		case "Name":
			emailData.Name = value
		case "Telephone":
			emailData.Phone = RemoveHTML(value)
		case "Email":
			emailData.EmailAddress = RemoveHTML(value)
		case "Personal status":
			emailData.Status = value
		case "Type of enquiry":
			emailData.EnquiryType = value
		}
	}

	lotRe := regexp.MustCompile(`colspan="2"><a href=.*?><img src="(.*?)"[\s\S]*?<tbody.*?><tr.*?>[\s]+<td.*?>[\s]*(.*?)</td>[\s\S]+?<tr.*>[\s]*<td.*?>[\s]*(.*?)</td>[\s\S]*?<tr.*>[\s]*<td.*?>[\s]*(.*?)</td>`)
	lotMatch := lotRe.FindStringSubmatch(body)

	if len(lotMatch) == 5 {
		emailData.ImageURL = lotMatch[1]
		emailData.Price = lotMatch[2]
		emailData.Description = lotMatch[3]
		emailData.Address = lotMatch[4]
	}

	linkRe := regexp.MustCompile(`<a.*href="(.+?)".*>View Property Details</a>`)
	linkMatch := linkRe.FindStringSubmatch(body)

	if len(linkMatch) == 2 {
		emailData.Link = linkMatch[1]
	}
}
