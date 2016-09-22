package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/text/language"
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

	t, e := template.ParseGlob("templates/*/*.tmpl")
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
	Lang            string
	ContactThanks   bool
	ContactSubjects map[string]string
}

func makeTabHandler(name string) func(w http.ResponseWriter, r *http.Request) {
	m := language.NewMatcher([]language.Tag{
		language.English,
		language.French})
	return func(w http.ResponseWriter, r *http.Request) {
		tags, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
		if err != nil {
			log.Printf("accept language error: %s\n", err)
		}
		tag, _, _ := m.Match(tags...)
		b, _ := tag.Base()
		code := b.ISO3()
		tName := fmt.Sprintf("%s/%s", code, name)
		err = theTemplate.ExecuteTemplate(w, tName, &TemplData{Active: name, Lang: code})
		if err != nil {
			log.Printf("error: %s", err)
		}
		logRequest(r)
	}
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

func letterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/pdf")
	fp, e := os.Open("letter.pdf")
	if e != nil {
		log.Printf("error: %s", e)
		return
	}
	defer fp.Close()
	_, e = io.Copy(w, fp)
	if e != nil {
		log.Printf("error: %s", e)
		return
	}
}

func main() {
	http.HandleFunc("/", makeTabHandler("root"))
	http.HandleFunc("/favicon.ico", favicoHandler)
	http.HandleFunc("/favicon", favicoHandler)
	http.HandleFunc("/mark.png", markHandler)
	http.HandleFunc("/mission", makeTabHandler("mission"))
	http.HandleFunc("/womb", makeTabHandler("womb"))
	http.HandleFunc("/about", makeTabHandler("about"))
	http.HandleFunc("/jobs", makeTabHandler("jobs"))
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/style.css", cssHandler)
	http.HandleFunc("/letter.pdf", letterHandler)

	log.Printf("serving from %s:%d in %s", *hostFlag, *portFlag, *rootDirFlag)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *hostFlag, *portFlag), nil)
	if err != nil {
		log.Fatalf("error %s", err)
	}
	log.Println("shutting down.")
}
