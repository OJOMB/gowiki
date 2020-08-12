package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gowiki/pages/page"
)

func handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the index page")
	}
}

func handleView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/view/"):]
		p, err := page.loadPage(title)
		if err != nil {
			fmt.Fprintf(w, "500 Internal Server Error whilst loading page")
		}
		fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
	}
}

func handleEdit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/edit/"):]
		p, err := page.loadPage(title)
		if err != nil {
			p = &page.Page{Title: title}
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
		p, err := page.loadPage(title)
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
