package twilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"olivia/models"
	"strings"
)

const BASE_URL = "https://api.twilio.com/2010-04-01/"

const TEST_SID = "AC2fea0226cd559f4f31ab3e4349887780"
const TEST_AUTHTOKEN = "79c4552a9679c670fc0dfe9e6eb9158d"

const LIVE_SID = "AC90ea5be1d88fbbb1479ad8c7460c9d1d"
const LIVE_AUTHTOKEN = "d9d29af856b1aabc49e37347721dd287"

const TESTING_MODE = false

type PhoneNumberCapabilities struct {
	Voice bool `json:"voice"`
	SMS   bool `json:"SMS"`
	MMS   bool `json:"MMS"`
	Fax   bool `json:"fax"`
}

type AvailablePhoneNumber struct {
	FriendlyName        string                  `json:"friendly_name"`
	PhoneNumber         string                  `json:"phone_number"`
	Lata                string                  `json:"lata"`
	RateCenter          string                  `json:"rate_center"`
	Latitude            string                  `json:"latitude"`
	Longitude           string                  `json:"longitude"`
	Region              string                  `json:"region"`
	ZIP                 string                  `json:"postal_code"`
	Country             string                  `json:"iso_country"`
	AddressRequirements string                  `json:"address_requirements"`
	Beta                bool                    `json:"beta"`
	Capabilities        PhoneNumberCapabilities `json:"capabilities"`
}

type PhoneNumber struct {
	SID         string `json:"string"`
	PhoneNumber string `json:"phone_number"`
}

type LookupResponse struct {
	CallerName     string `json:"caller_name"`
	CountryCode    string `json:"country_code"`
	PhoneNumber    string `json:"phone_number"`
	NationalFormat string `json:"national_format"`
}

func getCredentials() (string, string) {
	if TESTING_MODE {
		return TEST_SID, TEST_AUTHTOKEN
	}

	return LIVE_SID, LIVE_AUTHTOKEN
}

func GetPhoneNumbers(countryCode string) []AvailablePhoneNumber {
	var sid, authtoken = getCredentials()
	var client = &http.Client{}

	var URL string
	URL = fmt.Sprintf("%sAccounts/%s/AvailablePhoneNumbers/%s/Local.json", BASE_URL, sid, countryCode)

	req, reqErr := http.NewRequest("GET", URL, strings.NewReader(url.Values{
	//"AreaCode": {"510"},
	}.Encode()))
	if reqErr != nil {
		log.Println(reqErr.Error())
	}

	req.SetBasicAuth(sid, authtoken)

	resp, respErr := client.Do(req)

	if respErr != nil {
		log.Println(reqErr.Error())
	}

	bodyBytes, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr.Error())
	}

	resp.Body.Close()

	type Response struct {
		AvailablePhoneNumbers []AvailablePhoneNumber `json:"available_phone_numbers"`
	}

	var response Response

	json.Unmarshal(bodyBytes, &response)

	return response.AvailablePhoneNumbers
}

func ConnectPhoneNumber(user models.User, phoneNumber AvailablePhoneNumber) PhoneNumber {
	var sid, authtoken = getCredentials()
	var client = &http.Client{}

	var URL string
	URL = fmt.Sprintf("%sAccounts/%s/IncomingPhoneNumbers/Local.json", BASE_URL, sid)

	var friendlyUserName string
	friendlyUserName = strings.Join([]string{user.FirstName, user.SecondName}, "")

	req, reqErr := http.NewRequest("POST", URL, strings.NewReader(url.Values{
		"FriendlyName": {friendlyUserName},
		"PhoneNumber":  {phoneNumber.PhoneNumber},
		"SmsUrl":       {"http://api.getolivia.co/sms"},
		"SmsMethod":    {"POST"},
	}.Encode()))
	if reqErr != nil {
		log.Println(reqErr.Error())
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(sid, authtoken)

	resp, respErr := client.Do(req)

	if respErr != nil {
		log.Println(reqErr.Error())
	}

	bodyBytes, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr.Error())
	}

	resp.Body.Close()

	var response PhoneNumber
	json.Unmarshal(bodyBytes, &response)

	return response
}

func SendMessage(to string, from string, message string) {
	var sid, authtoken = getCredentials()
	var client = &http.Client{}

	var URL string
	URL = fmt.Sprintf("%sAccounts/%s/Messages.json", BASE_URL, sid)

	req, reqErr := http.NewRequest("POST", URL, strings.NewReader(url.Values{
		"To":   {to},
		"From": {from},
		"Body": {message},
	}.Encode()))
	if reqErr != nil {
		log.Printf("Error sending send sms request: %s", reqErr.Error())
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(sid, authtoken)

	resp, respErr := client.Do(req)
	if respErr != nil {
		log.Println(reqErr.Error())
	}

	_, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Printf("Error parsing send sms response: %s", readErr.Error())
	}

	resp.Body.Close()
}

func Lookup(number string) (LookupResponse) {
	var sid, authtoken = getCredentials()
	var client = &http.Client{}

	var URL string
	URL = fmt.Sprintf("https://lookups.twilio.com/v1/PhoneNumbers/%s?Type=carrier", number)

	req, reqErr := http.NewRequest("GET", URL, strings.NewReader(""))
	if reqErr != nil {
		log.Printf("Error sending lookup request: %s", reqErr.Error())
	}

	req.SetBasicAuth(sid, authtoken)

	resp, respErr := client.Do(req)
	if respErr != nil {
		log.Println(reqErr.Error())
	}

	bodyBytes, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Printf("Error parsing lookup response: %s", readErr.Error())
	}
	resp.Body.Close()

	var response LookupResponse
	json.Unmarshal(bodyBytes, &response)

	if response.CountryCode == "" {
		response.CountryCode = "US"
	}

	return response
}
