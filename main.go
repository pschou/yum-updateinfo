package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/jeanfric/goembed/countingwriter"
	"github.com/pschou/go-xmltree"
)

var version string

var (
	//fromStr  = flag.String("from", "nobody@example.com", "Set FROM value")
	relStr = flag.String("rel", "RedHat Compatible", "Set RELEASE value")
	//refStr   = flag.String("ref", "http://redhat.com", "Set REFERENCE value")
	confFile   = flag.String("conf", "packages.yml", "Config file")
	collection = flag.String("collection", "update-rpms", "Name for collection")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "YUM UpdateInfo Generator, Version", version, "(https://github.com/pschou/yum-updateinfo)")
		_, exec := path.Split(os.Args[0])
		fmt.Fprintln(os.Stderr, "Usage:\n  "+exec+" [options] path_to_repodata/repomd.xml\nOptions:")
		flag.PrintDefaults()
	}
	flag.Parse()

	var c Config
	c.getConf(*confFile)

	if flag.NArg() != 1 {
		log.Fatal("one argument required, repomd.xml file path")
	}

	// Get the path of the primary file from the repomd xml file
	primary := getPrimaryFile(flag.Arg(0))

	// Load the list of packages
	packages := loadPackages(primary)

	// Trim down the list to only the newest
	packages = findNewestPackages(packages)

	// Convert the package list to updates list
	var updates Updates
	for _, p := range packages {
		for _, t := range c.Packages {
			if t.reg.MatchString(p.Name) {
				updates.Update = append(updates.Update, makeUpdate(p, t.From, *relStr, t.RefTitle, t.RefURL))
			}
		}
	}

	// Generate the updateinfo xml file
	fileOutput := bytes.NewBuffer(nil)
	outerHash := sha256.New()
	outerWriter := io.MultiWriter(fileOutput, outerHash)

	gzipOutput := gzip.NewWriter(outerWriter)
	countingWriter := countingwriter.New(gzipOutput)
	innerHash := sha256.New()
	innerWriter := io.MultiWriter(countingWriter, innerHash)

	innerWriter.Write([]byte("<?xml version=\"1.0\"?>\n"))
	enc := xml.NewEncoder(innerWriter)
	enc.Indent("", "  ")
	enc.Encode(updates)
	enc.Flush()
	gzipOutput.Close()
	updatePath := fmt.Sprintf("repodata/%02x-updateinfo.xml.gz", outerHash.Sum(nil))
	updateLocal := path.Join(path.Dir(primary), fmt.Sprintf("%02x-updateinfo.xml.gz", outerHash.Sum(nil)))
	fmt.Println("writing to ", updateLocal)
	fh, err := os.Create(updateLocal)
	if err != nil {
		log.Fatal(err)
	}
	size, err := io.Copy(fh, fileOutput)
	if err != nil {
		log.Fatal(err)
	}

	updateData := xmltree.Element{
		StartElement: xml.StartElement{
			Name: xml.Name{Local: "data"},
			Attr: []xml.Attr{xml.Attr{Name: xml.Name{Local: "type"}, Value: "updateinfo"}},
		},
	}
	updateData.Children = []xmltree.Element{
		xmltree.Element{
			StartElement: xml.StartElement{
				Name: xml.Name{Local: "checksum"},
				Attr: []xml.Attr{xml.Attr{Name: xml.Name{Local: "type"}, Value: "sha256"}},
			},
			Content: fmt.Sprintf("%02x", outerHash.Sum(nil)),
		},
		xmltree.Element{
			StartElement: xml.StartElement{
				Name: xml.Name{Local: "open-checksum"},
				Attr: []xml.Attr{xml.Attr{Name: xml.Name{Local: "type"}, Value: "sha256"}},
			},
			Content: fmt.Sprintf("%02x", innerHash.Sum(nil)),
		},
		xmltree.Element{
			StartElement: xml.StartElement{
				Name: xml.Name{Local: "location"},
				Attr: []xml.Attr{xml.Attr{Name: xml.Name{Local: "href"}, Value: updatePath}},
			},
		},
		xmltree.Element{
			StartElement: xml.StartElement{
				Name: xml.Name{Local: "timestamp"},
			},
			Content: fmt.Sprintf("%d", time.Now().Unix()),
		},
		xmltree.Element{
			StartElement: xml.StartElement{
				Name: xml.Name{Local: "size"},
			},
			Content: fmt.Sprintf("%d", size),
		},
		xmltree.Element{
			StartElement: xml.StartElement{
				Name: xml.Name{Local: "open-size"},
			},
			Content: fmt.Sprintf("%d", countingWriter.BytesWritten()),
		},
	}

	// Insert the new xml file into the repomd xml file
	insert(flag.Arg(0), updateData)
}
