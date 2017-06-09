package zendesk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var _subdomain string
var _username string
var _password string

type TicketType string

const (
	Problem  TicketType = "problem"
	Incident TicketType = "incident"
	Question TicketType = "question"
	Task     TicketType = "task"
)

type TicketStatus string

const (
	New     TicketStatus = "new"
	Open    TicketStatus = "open"
	Pending TicketStatus = "pending"
	Hold    TicketStatus = "hold"
	Solved  TicketStatus = "solved"
	Closed  TicketStatus = "closed"
)

type TicketCommentType string

const (
	Comment      TicketCommentType = "Comment"
	VoiceComment TicketCommentType = "VoiceComment"
)

type TicketComment struct {
	Id   int64             `json:"id,omitempty"`
	Type TicketCommentType `json:"type,omitempty"`
	Body string            `json:"body,omitempty"`
}

type Ticket struct {
	Id          int64         `json:"id,omitempty"`
	Subject     string        `json:"subject,omitempty"`
	Comment     TicketComment `json:"comment,omitempty"`
	RequesterId int64         `json:"requester_id,omitempty"`
	SubmitterId int64         `json:"submitter_id,omitempty"`
	Type        TicketType    `json:"type,omitempty"`
	Status      TicketStatus  `json:"status,omitempty"`
}

func api(method, path string, params interface{}) []byte {
	trn := &http.Transport{}

	client := &http.Client{
		Transport: trn,
	}

	paramsBytes, jsonError := json.Marshal(params)
	if jsonError != nil {
		log.Println(jsonError)
	}

	req, reqErr := http.NewRequest(method, "https://"+_subdomain+".zendesk.com/api/v2/"+path, bytes.NewReader(paramsBytes))
	if reqErr != nil {
		log.Println(reqErr)
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(_username, _password)

	resp, doErr := client.Do(req)
	if doErr != nil {
		log.Println(doErr)
	}

	if method == "PUT" {
		log.Println(string(paramsBytes))
	}

	defer resp.Body.Close()

	respBytes, _ := ioutil.ReadAll(resp.Body)

	return respBytes
}

func SetCredentials(subdomain, username, password string) {
	_subdomain = subdomain
	_username = username
	_password = password
}

func CreateTicket(subject, content string) *Ticket {
	var ticket Ticket

	ticket.Subject = subject
	ticket.Comment.Body = content
	ticket.Status = New

	type Request struct {
		Ticket Ticket `json:"ticket"`
	}

	var request Request
	request.Ticket = ticket

	resp := api("POST", "tickets.json", request)

	type Response struct {
		Ticket Ticket `json:"Ticket"`
	}

	var response Response

	json.Unmarshal(resp, &response)

	return &response.Ticket
}

func GetTicket(id int64) *Ticket {
	resp := api("GET", fmt.Sprintf("tickets/%d.json", id), nil)

	type Response struct {
		Ticket Ticket `json:"Ticket"`
	}

	var response Response

	json.Unmarshal(resp, &response)

	return &response.Ticket
}

func (t *Ticket) Update() *Ticket {
	type Request struct {
		Ticket Ticket `json:"ticket"`
	}

	//var updatingTicket Ticket
	//updatingTicket.Comment = t.Comment
	//updatingTicket.Status = t.Status

	var request Request
	request.Ticket = *t

	log.Println(request)

	resp := api("PUT", fmt.Sprintf("tickets/%d.json", t.Id), request)

	type Response struct {
		Ticket Ticket `json:"Ticket"`
	}

	var response Response

	json.Unmarshal(resp, &response)
	//t = response.Ticket

	log.Println(string(resp))

	return &response.Ticket
}
