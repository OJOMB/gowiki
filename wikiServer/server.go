package wikiServer

import (
	"fmt"
	"gowiki/wikiLog"
	"gowiki/wikiPages"
	"gowiki/wikiRepo"
	"html/template"
	"net/http"
	"strings"
)

type Server struct {
	router *http.ServeMux
	repo   wikiRepo.Repo
	logger *wikiLog.Logger
}

func GetWikiServer(repo wikiRepo.Repo, logger *wikiLog.Logger) *Server {
	mux := http.NewServeMux()
	var s *Server = &Server{
		router: mux,
		repo:   repo,
		logger: logger,
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	// register handlers against patterns in the DefaultServeMux
	s.router.HandleFunc("/", s.HandleIndex())
	s.router.HandleFunc("/edit/", s.HandleEdit())
	s.router.HandleFunc("/save/", s.HandleSave())
	s.router.HandleFunc("/view/", s.HandleView())
}

func (s *Server) HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the index page bro")
	}
}

func (s *Server) HandleView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/view/"):]
		p, err := s.repo.ReadPage(title)
		if err != nil {
			if cause := err.Error(); strings.Contains(cause, "no such file or directory") {
				renderTemplate(w, "page_not_found", &wikiPages.Page{Title: title})
				return
			}
			fmt.Fprintf(w, "500 Internal Server Error whilst loading page")
			return
		}
		renderTemplate(w, "view", p)
	}
}

func (s *Server) HandleEdit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/edit/"):]
		p, err := s.repo.ReadPage(title)
		if err != nil {
			p = &wikiPages.Page{Title: title}
		}
		renderTemplate(w, "edit", p)
	}
}

func (s *Server) HandleSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[len("/save/"):]
		body := r.FormValue("body")
		p := &wikiPages.Page{Title: title, Body: []byte(body)}
		s.repo.WritePage(p)
		http.Redirect(w, r, "/view/"+title, http.StatusFound)
	}
}

func renderTemplate(w http.ResponseWriter, tName string, p *wikiPages.Page) {
	t, err := template.ParseFiles("wikiTemplates/" + tName + ".html")
	if err != nil {
		fmt.Fprintln(w, "Failed to render edit page template")
	}
	t.Execute(w, p)
}
