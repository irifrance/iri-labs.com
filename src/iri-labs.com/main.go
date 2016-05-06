package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var favIcon []byte = nil
var markPng []byte = nil
var styleCss []byte = nil
var theTemplate *template.Template = nil

func init() {
	fi, e := readFile("data/favicon.ico")
	if e != nil {
		log.Println(e)
		return
	}
	favIcon = fi

	m, e := readFile("data/mark.png")
	if e != nil {
		log.Println(e)
		return
	}
	markPng = m

	css, e := readFile("style.css")
	if e != nil {
		log.Println(e)
		return
	}
	styleCss = css

	t, e := template.ParseGlob("templates/*.tmpl")
	if e != nil {
		log.Println(e)
		return
	}
	theTemplate = t
}

func readFile(p string) ([]byte, error) {
	f, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	st, e := f.Stat()
	if e != nil {
		return nil, e
	}
	dat := make([]byte, st.Size())
	n, e := f.Read(dat)
	if int64(n) != st.Size() {
		return nil, e
	}
	return dat, nil
}

func logRequest(r *http.Request) {
	log.Printf("%s %s %s %s %s", r.RemoteAddr, r.Host, r.RequestURI, r.URL, r.Method)
}

type TemplData struct {
	Active string
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	err := theTemplate.ExecuteTemplate(w, "root", &TemplData{Active: "/"})
	if err != nil {
		log.Printf("error: %s", err)
	}
	logRequest(r)
}

func missionHandler(w http.ResponseWriter, r *http.Request) {
	err := theTemplate.ExecuteTemplate(w, "mission", &TemplData{Active: "mission"})
	if err != nil {
		log.Printf("error: %s", err)
	}
	logRequest(r)
}

func wombHandler(w http.ResponseWriter, r *http.Request) {
	err := theTemplate.ExecuteTemplate(w, "womb", &TemplData{Active: "womb"})
	if err != nil {
		log.Printf("error: %s", err)
	}
	logRequest(r)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	err := theTemplate.ExecuteTemplate(w, "about", &TemplData{Active: "about"})
	if err != nil {
		log.Printf("error: %s", err)
	}
	logRequest(r)
}

func jobsHandler(w http.ResponseWriter, r *http.Request) {
	err := theTemplate.ExecuteTemplate(w, "jobs", &TemplData{Active: "jobs"})
	if err != nil {
		log.Printf("error: %s", err)
	}
	logRequest(r)
}

func favicoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(favIcon)
}

func markHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(markPng)
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.Write(styleCss)
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/favicon.ico", favicoHandler)
	http.HandleFunc("/favicon", favicoHandler)
	http.HandleFunc("/mark.png", markHandler)
	http.HandleFunc("/mission", missionHandler)
	http.HandleFunc("/womb", wombHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/jobs", jobsHandler)
	http.HandleFunc("/style.css", cssHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("error %s", err)
	}
	log.Println("shutting down.")
}
