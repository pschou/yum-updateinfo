package main

import (
	"encoding/xml"
	"strings"
	"time"

	"github.com/araddon/dateparse"

	rpmver "github.com/pschou/go-rpm-version"
)

type Updates struct {
	XMLName xml.Name `xml:"updates"`
	Text    string   `xml:",chardata"`
	Update  []Update `xml:"update"`
}

type Update struct {
	Text    string `xml:",chardata"`
	From    string `xml:"from,attr"`
	Status  string `xml:"status,attr"`
	Type    string `xml:"type,attr"`
	Version string `xml:"version,attr"`
	ID      string `xml:"id"`
	Title   string `xml:"title"`
	Issued  struct {
		Text string `xml:",chardata"`
		Date string `xml:"date,attr"`
	} `xml:"issued"`
	Updated struct {
		Text string `xml:",chardata"`
		Date string `xml:"date,attr"`
	} `xml:"updated"`
	Rights      string `xml:"rights"`
	Release     string `xml:"release"`
	PushCount   string `xml:"pushcount"`
	Severity    string `xml:"severity"`
	Summary     string `xml:"summary"`
	Description string `xml:"description"`
	Solution    string `xml:"solution"`
	References  struct {
		Text      string `xml:",chardata"`
		Reference struct {
			Text  string `xml:",chardata"`
			Href  string `xml:"href,attr"`
			ID    string `xml:"id,attr"`
			Type  string `xml:"type,attr"`
			Title string `xml:"title,attr"`
		} `xml:"reference"`
	} `xml:"references"`
	Pkglist struct {
		Text       string `xml:",chardata"`
		Collection struct {
			Text    string `xml:",chardata"`
			Short   string `xml:"short,attr"`
			Name    string `xml:"name"`
			Package struct {
				Text     string `xml:",chardata"`
				Name     string `xml:"name,attr"`
				Version  string `xml:"version,attr"`
				Release  string `xml:"release,attr"`
				Epoch    int    `xml:"epoch,attr"`
				Arch     string `xml:"arch,attr"`
				Filename string `xml:"filename"`
				Sum      struct {
					Text string `xml:",chardata"`
					Type string `xml:"type,attr"`
				} `xml:"sum"`
			} `xml:"package"`
		} `xml:"collection"`
	} `xml:"pkglist"`
}

func makeUpdate(p Package, from, release, reference_title, reference_url string) Update {
	ver := rpmver.NewVersionWithValues(p.Version.Ver, p.Version.Rel, p.Version.Epoch)
	ret := Update{
		Status:      "stable",
		Type:        "security",
		Version:     "1",
		From:        from,
		ID:          p.Name + "-" + ver.String() + "." + p.Arch,
		Title:       p.Name + " update",
		Rights:      strings.TrimSpace("Copyright " + p.Format.Vendor),
		Release:     release,
		PushCount:   "1",
		Severity:    "None",
		Solution:    "Updates for fixing bugs or enhancement of " + p.Name,
		Summary:     "New " + p.Name + " packages are available",
		Description: p.Description,
	}
	if val, err := dateparse.ParseAny(p.Time.Build); err == nil {
		ret.Issued.Date = val.Add(2 * time.Second).Format("2006-01-02 15:04:05")
		ret.Updated.Date = val.Format("2006-01-02 15:04:05")
	} else {
		ret.Issued.Date = p.Time.Build
		ret.Updated.Date = ret.Issued.Date
	}
	ret.References.Reference.Href = reference_url
	ret.References.Reference.Type = "self"
	ret.References.Reference.ID = ret.ID
	ret.References.Reference.Title = reference_title
	ret.Pkglist.Collection.Short = *collection
	ret.Pkglist.Collection.Name = *collection
	ret.Pkglist.Collection.Package.Name = p.Name
	ret.Pkglist.Collection.Package.Version = p.Version.Ver
	ret.Pkglist.Collection.Package.Release = p.Version.Rel
	ret.Pkglist.Collection.Package.Epoch = p.Version.Epoch
	ret.Pkglist.Collection.Package.Arch = p.Arch
	ret.Pkglist.Collection.Package.Filename = p.Location.Href
	ret.Pkglist.Collection.Package.Sum.Type = p.Checksum.Type
	ret.Pkglist.Collection.Package.Sum.Text = p.Checksum.Text
	return ret
}
