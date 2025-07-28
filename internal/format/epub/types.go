package epub

import "encoding/xml"

// Internal XML structures for EPUB parsing
type container struct {
	XMLName   xml.Name   `xml:"container"`
	RootFiles []rootFile `xml:"rootfiles>rootfile"`
}

type rootFile struct {
	FullPath  string `xml:"full-path,attr"`
	MediaType string `xml:"media-type,attr"`
}

type Package struct {
	XMLName  xml.Name `xml:"package"`
	Metadata Metadata `xml:"metadata"`
	Manifest Manifest `xml:"manifest"`
	Spine    Spine    `xml:"spine"`
}

type Metadata struct {
	Title     string `xml:"title"`
	Creator   string `xml:"creator"`
	Publisher string `xml:"publisher"`
	Language  string `xml:"language"`
}

type Manifest struct {
	Items []Item `xml:"item"`
}

type Item struct {
	ID        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr"`
}

type Spine struct {
	ItemRefs []ItemRef `xml:"itemref"`
}

type ItemRef struct {
	IDRef string `xml:"idref,attr"`
}
