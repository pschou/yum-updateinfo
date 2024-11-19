package main

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
)

type Repomd struct {
	XMLName  xml.Name     `xml:"repomd"`
	Text     string       `xml:",chardata"`
	Xmlns    string       `xml:"xmlns,attr"`
	Rpm      string       `xml:"rpm,attr"`
	Revision string       `xml:"revision"`
	Data     []RepomdData `xml:"data"`
}
type RepomdData struct {
	Text     string `xml:",chardata"`
	Type     string `xml:"type,attr"`
	Checksum struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
	} `xml:"checksum"`
	OpenChecksum struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
	} `xml:"open-checksum"`
	Location struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
	} `xml:"location"`
	Timestamp string `xml:"timestamp"`
	Size      string `xml:"size"`
	OpenSize  string `xml:"open-size"`
}

func getPrimaryFile(repomdFile string) string {
	repomd := loadRepomd(repomdFile)
	for _, e := range repomd.Data {
		if e.Type == "primary" {
			locationDir := path.Dir(e.Location.Href)
			repomdPath := path.Dir(repomdFile)
			parentDir, repomdDir := path.Split(repomdPath)
			if strings.EqualFold(repomdDir, locationDir) {
				// Found path
			} else if strings.EqualFold(locationDir, "repodata") {
				// Fall back for misnamed paths
				fmt.Printf("Warning: repomd file's parent directory %q doesn't match the parent directory for primary xml %q.\n", repomdDir, locationDir)
			} else {
				log.Fatal("repomdFile must be in the same path as the primary xml file.")
			}

			// Join the parts together
			joined, err := url.JoinPath(parentDir, e.Location.Href)
			if err != nil {
				log.Fatal(err)
			}
			return joined
		}
	}
	log.Fatal("Could not find \"primary\" type in repomd.xml")
	return ""
}

func loadRepomd(file string) Repomd {
	// Open the repomd file
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
	var ret Repomd
	err = xmlDecoder.Decode(&ret)
	if err != nil {
		log.Fatal(err)
	}

	return ret
}
