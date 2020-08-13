package wikiRepo

import "gowiki/wikiPages"

type Repo interface {
	ReadPage(title string) (*wikiPages.Page, error)
	WritePage(p *wikiPages.Page) error
}
