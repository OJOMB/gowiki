package wikiPages

// a wiki is a series of interconnected pages that can be modeled
// as being composed of titles and bodies of content.
type Page struct {
	Title string
	Body  []byte // byte slice rather than string because the io libs expect it in line with HTTP
}
