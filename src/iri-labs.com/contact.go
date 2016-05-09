package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sync"
)

var contactSubjects = map[string]string{"jobs": "Jobs",
	"info":      "General Info",
	"available": "Request availability",
	"triage":    "Request triage",
	"join":      "Join our network"}

type Contact struct {
	NetAddr   string
	Name      string
	Institute string
	Phone     string
	Country   string
	Subject   string
	Body      string
}

var contactLog string = "log/contact"

var contactKeys = []string{
	"name", "institute", "phone", "country", "subject", "body"}

var contactLogMutex sync.Mutex

func contactHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if r.Method == "POST" {
		if e := r.ParseForm(); e != nil {
			log.Printf("contact error: %s", e)
		}
		c := &Contact{}
		for _, k := range contactKeys {
			v := r.FormValue(k)
			switch k {
			case "name":
				c.Name = v
			case "institute":
				c.Institute = v
			case "phone":
				c.Phone = v
			case "country":
				c.Country = v
			case "subject":
				c.Subject = v
			case "body":
				c.Body = v
			default:
				panic("unknown key")
			}
		}
		if e := logContact(c); e != nil {
			log.Printf("error log Contact: %s", e)
		}
		if e := sendContact(c); e != nil {
			log.Printf("error send Contact: %s", e)
		}
	}
	// GET or Post
	err := theTemplate.ExecuteTemplate(w, "contact", &TemplData{Active: "jobs", ContactThanks: r.Method == "POST"})
	if err != nil {
		log.Printf("error: %s", err)
	}
}

func logContact(c *Contact) error {
	contactLogMutex.Lock()
	defer contactLogMutex.Unlock()
	fp, e := os.OpenFile(contactLog, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if e != nil {
		return e
	}
	defer fp.Close()
	enc := json.NewEncoder(fp)
	if e := enc.Encode(c); e != nil {
		return e
	}
	return nil
}

func sendContact(c *Contact) error {
	b := bytes.NewBuffer(nil)
	srvAddr := "aspmx.l.google.com:25"
	sender := "www@iri-labs.com"
	rcpts := []string{"wsc@iri-labs.com"}
	msg := `From: www@iri-labs.com
To: wsc@iri-labs.com
Subject: new contact request

`
	b.WriteString(msg)
	enc := json.NewEncoder(b)
	if e := enc.Encode(c); e != nil {
		return e
	}
	if e := smtp.SendMail(srvAddr, nil, sender, rcpts, b.Bytes()); e != nil {
		return e
	}
	return nil
}
