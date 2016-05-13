package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var favIcon []byte = nil
var markPng []byte = nil
var styleCss []byte = nil
var theTemplate *template.Template = nil
var rootDirFlag *string = flag.String("root", ".", "root directory from which to serve")
var logDirFlag *string = flag.String("log", "", "log directory")
var hostFlag *string = flag.String("host", "", "default host address on which to serve")
var portFlag *int = flag.Int("port", 80, "serve on this port")

func init() {
	flag.Parse()
	if e := os.Chdir(*rootDirFlag); e != nil {
		log.Fatal("unable to chdir to '%s'", *rootDirFlag)
	}

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
	Active          string
	ContactThanks   bool
	ContactSubjects map[string]string
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
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/style.css", cssHandler)

	log.Printf("serving from %s:%d in %s", *hostFlag, *portFlag, *rootDirFlag)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *hostFlag, *portFlag), nil)
	if err != nil {
		log.Fatalf("error %s", err)
	}
	log.Println("shutting down.")
}
