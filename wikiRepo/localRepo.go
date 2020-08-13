package wikiRepo

import (
	"fmt"
	"gowiki/wikiLog"
	"gowiki/wikiPages"
	"gowiki/wikiUtils"
	"io/ioutil"

	"github.com/pkg/errors"
)

const PAGE_SAVE_FILE_EXTENSION = ".txt"

type LocalRepo struct {
	filepath []string
	logger   *wikiLog.Logger
}

func GetLocalWikiRepo(logger *wikiLog.Logger) *LocalRepo {
	return &LocalRepo{
		filepath: []string{"wikiPages", "dataFiles"},
		logger:   logger,
	}
}

// Writes Page data to txt files stored locally
func (repo *LocalRepo) WritePage(p *wikiPages.Page) error {
	filePath := repo.getPageDataFilePath(p.Title)
	// ioutil.WriteFile writes a byte slice to a file and returns an error
	return errors.Wrap(
		ioutil.WriteFile(
			filePath,
			p.Body,
			0600, // rw permissions for current user only
		),
		fmt.Sprintf("failed to write Page struct %s to file", p.Title),
	)
}

// Loads Page data from txt files stored locally
func (repo *LocalRepo) ReadPage(title string) (*wikiPages.Page, error) {
	filePath := repo.getPageDataFilePath(title)
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to load page data from file %s", filePath))
	}
	return &wikiPages.Page{Title: title, Body: body}, nil
}

func (repo *LocalRepo) getPageDataFilePath(title string) string {
	s := append(repo.filepath, title)
	return wikiUtils.ConstructPath(s, PAGE_SAVE_FILE_EXTENSION)
}
