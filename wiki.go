package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

const SAVE_FILE_FORMAT = "txt"

// a wiki is a series of interconnected pages that can be modeled
// as being composed of titles and bodies of content.
type Page struct {
	Title string
	Body  []byte // byte slice rather than string because the io libs expect it in line with HTTP
}

// the Page struct gives us a model for Page related data in memory
// for persistent storage we need a method that can save the data to a file
func (p *Page) save() error {
	filename := "pages" + string(os.PathSeparator) + p.Title + "." + SAVE_FILE_FORMAT
	// ioutil.WriteFile writes a byte slice to a file and returns an error
	return errors.Wrap(
		ioutil.WriteFile(
			filename,
			p.Body,
			0600, // rw permissions for current user only
		),
		fmt.Sprintf("failed to write Page struct %s to file", p.Title),
	)
}

// since we're saving pages to files, we need a method to load them back again when needed
func loadPage(title string) (*Page, error) {
	filename := "pages" + string(os.PathSeparator) + title + "." + SAVE_FILE_FORMAT
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to load page from file %s.txt", filename))
	}
	return &Page{title, body}, nil
}

func handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the index page")
	}
}

func handleView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/view/"):]
		p, err := loadPage(title)
		if err != nil {
			fmt.Fprintf(w, "500 Internal Server Error whilst loading page")
		}
		fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
	}
}

func handleEdit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/edit/"):]
		p, err := loadPage(title)
		if err != nil {
			p = &Page{Title: title}
		}
		fmt.Fprintf(
			w,
			"<h1>Editing %s</h1>"+
				"<form action=\"/save/%s\" method=\"POST\">"+
				"<textarea name=\"body\">%s</textarea><br>"+
				"<input type=\"submit\" value=\"Save\">"+
				"</form>",
			p.Title,
			p.Title,
			p.Body,
		)
	}
}

func handleSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/view/"):]
		p, err := loadPage(title)
		if err != nil {
			fmt.Fprintf(w, "500 Internal Server Error whilst loading data for page")
		}

		t, _ := template.ParseFiles("/templates/edit.html")
		if err != nil {
			fmt.Fprintln(w, "Failed to render edit page template")
		}
		t.Execute(w, p)
	}
}

func main() {
	// HandleFunc registers a function against a pattern in the DefaultServeMux
	http.HandleFunc("/", handleIndex())
	http.HandleFunc("/edit/", handleEdit())
	http.HandleFunc("/save/", handleSave())
	http.HandleFunc("/view/", handleView())

	// ListenAndServe only returns in the event of an exceptional error
	// It's therefore good practice to log it's return value like so:
	log.Fatal(http.ListenAndServe(":8080", nil))

}
