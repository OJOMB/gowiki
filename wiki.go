package main

import (
	"log"
	"net/http"

	"gowiki/wikiLog"
	"gowiki/wikiRepo"
	"gowiki/wikiServer"
)

func main() {
	var logger *wikiLog.Logger = wikiLog.GetWikiLogger()
	var s *wikiServer.Server = wikiServer.GetWikiServer(
		wikiRepo.GetLocalWikiRepo(logger),
		logger,
	)

	// ListenAndServe only returns in the event of an exceptional error
	// It's therefore good practice to log it's return value like so:
	log.Fatal(http.ListenAndServe(":8080", s))

}
