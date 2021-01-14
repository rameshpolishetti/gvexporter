package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type GlobalVariable struct {
	XMLName xml.Name `xml:"globalVariable"`
	Name    string   `xml:"name"`
	Value   string   `xml:"value"`
}

type GlobalVariables struct {
	XMLName        xml.Name         `xml:"globalVariables"`
	GlobalVariable []GlobalVariable `xml:"globalVariable"`
}

type Repository struct {
	XMLName         xml.Name        `xml:"repository"`
	GlobalVariables GlobalVariables `xml:"globalVariables"`
}

func main() {

	beProjPath := flag.String("i", "", "Input BE Project Path")
	outputFilename := flag.String("o", "", "Output Filename")
	flag.Parse()

	if *beProjPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *outputFilename == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Exporting GV definitions from the BE Project[%s] into the file[%s]\n", *beProjPath, *outputFilename)

	// read GVs into map
	gvMap := make(map[string]string)
	gvsPath := filepath.Join(*beProjPath, "defaultVars")
	err := filepath.Walk(gvsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if strings.HasSuffix(path, "defaultVars.substvar") {
			keyPrefix := strings.TrimPrefix(path, gvsPath+"/")
			keyPrefix = strings.TrimSuffix(keyPrefix, "defaultVars.substvar")

			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			bytes, _ := ioutil.ReadAll(file)
			var xmlData Repository
			xml.Unmarshal(bytes, &xmlData)

			for _, gv := range xmlData.GlobalVariables.GlobalVariable {
				key := keyPrefix + gv.Name
				gvMap[key] = gv.Value
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// write gvs into output file
	outFile, err := os.OpenFile(*outputFilename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	gvJSON, err := json.MarshalIndent(gvMap, "", "  ")
	outFile.WriteString(string(gvJSON))

	fmt.Println("DONE")
}
