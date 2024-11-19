package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pschou/go-xmltree"
)

func insert(file string, updates xmltree.Element) {
	infile, err := os.OpenFile(file, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	log.Println("Successfully Opened " + file)
	// defer the closing of our xmlFile so that we can parse it later on

	root, err := xmltree.ParseXML(infile)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("doc: %#v\n", root)
	//os.Exit(0)
	//trimSpace(root)

	/*
		dbVersion := root.FindOne(&xmltree.Selector{
			Name: xml.Name{Local: "database_version"},
		})
		if dbVersion != nil {
			ver := *dbVersion
			ver.Name = xml.Name{Local: "database_version"}
			ver.Scope = xmltree.Scope{}
			updates.Children = append(updates.Children, ver)
		}
	*/

	data := root.FindOne(&xmltree.Selector{Depth: 1,
		Name: xml.Name{Local: "data"},
		Attr: []xml.Attr{xml.Attr{Name: xml.Name{Local: "type"}, Value: "updateinfo"}},
	})
	if data != nil {
		fmt.Println("found updateinfo, replacing")
		*data = updates
	} else {
		root.Children = append(root.Children, updates)
	}
	infile.Seek(0, io.SeekStart)
	log.Println("Writing output")
	infile.Write([]byte(xml.Header))
	err = xmltree.EncodeIndent(infile, root, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	loc, _ := infile.Seek(0, io.SeekCurrent)
	infile.Truncate(loc)
	log.Println("done, writing to file")

}
