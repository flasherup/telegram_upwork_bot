package upwork

import "encoding/xml"

type UWFeed struct {
	XMLName   xml.Name  `xml:"feed"`
	Xmlns     string    `xml:"xmlns,attr"`
	Id        string    `xml:"id"`
	Title     string    `xml:"title"`
	Author    UWAuthor  `xml:"author"`
	Updated   string    `xml:"updated"`
	Link      UWLink    `xml:"link"`
	Subtitle  string    `xml:"subtitle"`
	Rights    string    `xml:"rights"`
	Logo      string    `xml:"logo"`
	Generator string    `xml:"generator"`
	Entries   []UWEntry `xml:"entry"`
}

type UWAuthor struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name"`
	Email   string   `xml:"email"`
}

type UWLink struct {
	XMLName  xml.Name `xml:"link"`
	Rel      string   `xml:"rel,attr"`
	Href     string   `xml:"href,attr"`
	Hreflang string   `xml:"hreflang,attr"`
}

type UWEntry struct {
	XMLName xml.Name `xml:"entry"`
	Id      string   `xml:"id"`
	Title   string   `xml:"title"`
	Updated string   `xml:"updated"`
	Link    UWLink   `xml:"link"`
	Summary string   `xml:"summary"`
	Content string   `xml:"content"`
}
