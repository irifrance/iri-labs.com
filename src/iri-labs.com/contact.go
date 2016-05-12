package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"sync"
	"crypto/tls"
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
	err := theTemplate.ExecuteTemplate(w, "contact", &TemplData{Active: "contact",
		ContactThanks:   r.Method == "POST",
		ContactSubjects: contactSubjects})
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

// we do this manually due to constraints of gcloud networking on SMTP
// we need to set HELO host manually
// we need to use TLS
// smtp.SendMail doesn't do this, so we go one step lower
func sendContact(ctc *Contact) error {
	srvAddr := "smtp-relay.gmail.com:587"
	nc, err := net.Dial("tcp", srvAddr)
	if err != nil {
		return fmt.Errorf("net dial: %s", err)
	}
	c, err := smtp.NewClient(nc, "smtp-relay.gmail.com")
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	c.Hello("iri-labs.com")
	if err := c.StartTLS(&tls.Config{ServerName: "smtp-relay.gmail.com"}); err != nil {
		return fmt.Errorf("start-tls: %s", err)
	}
	if err := c.Mail("www@iri-labs.com"); err != nil {
		return fmt.Errorf("mail from: %s", err)
	}
	if err := c.Rcpt("wsc@iri-labs.com"); err != nil {
		return fmt.Errorf("rcpt err: %s", err)
	}
	bw, err := c.Data()
	if err != nil {
		return fmt.Errorf("data err: %s", err)
	}
	defer bw.Close()

	b := bytes.NewBufferString(`From: www@iri-labs.com
To: wsc@iri-labs.com
Subject: [www.iri-labs.com] new contact request

`)
	_, e := b.WriteTo(bw)
	if e != nil {
		return fmt.Errorf("write headers: %s", e)
	}
	b.Reset()

	jb, e := json.MarshalIndent(ctc, "", "\t")
	if e != nil {
		return fmt.Errorf("marshal: %s", e)
	}
	_, e = bw.Write(jb)
	if e != nil {
		return fmt.Errorf("write json: %s", e)
	}
	if e := c.Quit(); e != nil {
		return fmt.Errorf("quit: %s", e)
	}
	return nil
}
	
	

