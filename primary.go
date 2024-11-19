package main

import (
	"compress/gzip"
	"encoding/xml"
	"io"
	"log"
	"os"
	"strings"

	rpmver "github.com/pschou/go-rpm-version"
)

type PackageMetadata struct {
	XMLName  xml.Name  `xml:"metadata"`
	Text     string    `xml:",chardata"`
	Xmlns    string    `xml:"xmlns,attr"`
	Rpm      string    `xml:"rpm,attr"`
	Packages string    `xml:"packages,attr"`
	Package  []Package `xml:"package"`
}

type Package struct {
	Text    string `xml:",chardata"`
	Type    string `xml:"type,attr"`
	Name    string `xml:"name"`
	Arch    string `xml:"arch"`
	Version struct {
		Text  string `xml:",chardata"`
		Epoch int    `xml:"epoch,attr"`
		Ver   string `xml:"ver,attr"`
		Rel   string `xml:"rel,attr"`
	} `xml:"version"`
	Checksum struct {
		Text  string `xml:",chardata"`
		Type  string `xml:"type,attr"`
		Pkgid string `xml:"pkgid,attr"`
	} `xml:"checksum"`
	Summary     string `xml:"summary"`
	Description string `xml:"description"`
	Packager    string `xml:"packager"`
	URL         string `xml:"url"`
	Time        struct {
		Text  string `xml:",chardata"`
		File  string `xml:"file,attr"`
		Build string `xml:"build,attr"`
	} `xml:"time"`
	Size struct {
		Text      string `xml:",chardata"`
		Package   string `xml:"package,attr"`
		Installed string `xml:"installed,attr"`
		Archive   string `xml:"archive,attr"`
	} `xml:"size"`
	Location struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
	} `xml:"location"`
	Format struct {
		Text        string `xml:",chardata"`
		License     string `xml:"license"`
		Vendor      string `xml:"vendor"`
		Group       string `xml:"group"`
		Buildhost   string `xml:"buildhost"`
		Sourcerpm   string `xml:"sourcerpm"`
		HeaderRange struct {
			Text  string `xml:",chardata"`
			Start string `xml:"start,attr"`
			End   string `xml:"end,attr"`
		} `xml:"header-range"`
		Provides struct {
			Text  string `xml:",chardata"`
			Entry []struct {
				Text  string `xml:",chardata"`
				Name  string `xml:"name,attr"`
				Flags string `xml:"flags,attr"`
				Epoch string `xml:"epoch,attr"`
				Ver   string `xml:"ver,attr"`
				Rel   string `xml:"rel,attr"`
			} `xml:"entry"`
		} `xml:"provides"`
		Requires struct {
			Text  string `xml:",chardata"`
			Entry []struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"entry"`
		} `xml:"requires"`
		Conflicts struct {
			Text  string `xml:",chardata"`
			Entry struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"entry"`
		} `xml:"conflicts"`
		Obsoletes struct {
			Text  string `xml:",chardata"`
			Entry struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"entry"`
		} `xml:"obsoletes"`
		File []struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		} `xml:"file"`
	} `xml:"format"`
}

func findNewestPackages(in []Package) (out []Package) {
newestLoop:
	for _, e := range in {
		if e.Type != "rpm" {
			continue // Ignore non-rpm files
		}
		for i, t := range out {
			if t.Name == e.Name && strings.EqualFold(t.Arch, e.Arch) {
				tVer := rpmver.NewVersionWithValues(t.Version.Ver, t.Version.Rel, t.Version.Epoch)
				eVer := rpmver.NewVersionWithValues(e.Version.Ver, e.Version.Rel, e.Version.Epoch)
				if tVer.LessThan(eVer) {
					out[i] = e
				}
				continue newestLoop
			}
		}
		out = append(out, e)
	}
	return
}

func loadPackages(file string) []Package {
	// Open the package file
	xmlFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	var fileContents io.Reader

	// Test if gzipped
	gz, err := gzip.NewReader(xmlFile)
	if err == nil {
		fileContents = gz
		defer gz.Close()
	} else {
		xmlFile.Seek(0, io.SeekStart)
		fileContents = xmlFile
	}

	xmlDecoder := xml.NewDecoder(fileContents)
	var ret PackageMetadata
	err = xmlDecoder.Decode(&ret)
	if err != nil {
		log.Fatal(err)
	}

	return ret.Package
}
