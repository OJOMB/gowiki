package page

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const SAVE_FILE_FORMAT = "txt"
const SEPARATOR = string(os.PathSeparator)

// a wiki is a series of interconnected pages that can be modeled
// as being composed of titles and bodies of content.
type Page struct {
	Title string
	Body  []byte // byte slice rather than string because the io libs expect it in line with HTTP
}

// the Page struct gives us a model for Page related data in memory
// for persistent storage we need a method that can save the data to a file
func (p *Page) Save() error {
	filePath := getPageDataFilePath(p.Title)
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

// since we're saving pages to files, we need a method to load them back again when needed
func LoadPage(title string) (*Page, error) {
	filePath := getPageDataFilePath(title)
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to load page data from file %s.txt", filename))
	}
	return &Page{title, body}, nil
}

func getPageDataFilePath(title string) string {
	s := append([]string{"pages", "data"}, title)
	return constructPath(s, SAVE_FILE_FORMAT)
}

func constructPath(fields []string, fileExtension string) string {
	return strings.Join(fields, SEPARATOR) + fileExtension
}
