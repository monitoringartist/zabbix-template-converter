package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	v40 "github.com/monitoringartist/go-zabbix/v40"
	v44 "github.com/monitoringartist/go-zabbix/v44"
	"gopkg.in/yaml.v3"
)

type ZabbixExport struct {
	XMLName xml.Name `xml:"zabbix_export" json:"-" yaml:"-"`
	Version string   `xml:"version,omitempty" json:"version,omitempty" yaml:"version,omitempty"`
}

func main() {
	in := flag.String("in", "", "input template file with xml/json/yaml extension")
	out := flag.String("out", "stdout", "output template file with xml/json/yaml extension")
	flag.Parse()
	if *in == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *out == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	supportedExtentsion := map[string]string{
		"xml":  "",
		"json": "",
		"yaml": "",
	}

	sIn := strings.Split(string([]byte(*in)), ".")
	_, ok := supportedExtentsion[sIn[len(sIn)-1]]
	if !ok {
		fmt.Printf("Unsupported input format: %s\n", sIn[len(sIn)-1])
		os.Exit(1)
	}
	sOut := strings.Split(string([]byte(*out)), ".")
	_, ok = supportedExtentsion[sOut[len(sOut)-1]]
	if !ok {
		fmt.Printf("Unsupported output format: %s\n", sOut[len(sOut)-1])
		os.Exit(1)
	}

	inData, err := ioutil.ReadFile(*in)
	if err != nil {
		fmt.Printf("Can't open input file, error: %s\n", err.Error())
		os.Exit(1)
	}

	// version parsing
	templVersion := &ZabbixExport{}
	switch sIn[len(sIn)-1] {
	case "xml":
		err = xml.Unmarshal([]byte(inData), &templVersion)
		if err != nil {
			fmt.Printf("Input XML file unmarshal error: %s\n", err.Error())
			os.Exit(1)
		}
	case "json":
		err = json.Unmarshal([]byte(inData), &templVersion)
		if err != nil {
			fmt.Printf("Input JSON file unmarshal error: %s\n", err.Error())
			os.Exit(1)
		}
	case "yaml":
		err = yaml.Unmarshal([]byte(inData), &templVersion)
		if err != nil {
			fmt.Printf("Input YAML file unmarshal error: %s\n", err.Error())
			os.Exit(1)
		}
	}

	if templVersion.Version != "4.4" && templVersion.Version != "4.2" && templVersion.Version != "4.0" {
		fmt.Printf("Unsupported Zabbix template version %s, only versions 4.4/4.2/4.0 are supported\n", templVersion.Version)
		os.Exit(1)
	}

	var templ interface{}
	if templVersion.Version == "4.4" {
		//v4.4
		templ = &v44.ZabbixExport{}
	} else {
		// probably v4.2 = v4.0
		templ = &v40.ZabbixExport{}
	}

	switch sIn[len(sIn)-1] {
	case "xml":
		err = xml.Unmarshal([]byte(inData), &templ)
		if err != nil {
			fmt.Printf("Input XML file unmarshal error: %s\n", err.Error())
			os.Exit(1)
		}
	case "json":
		err = json.Unmarshal([]byte(inData), &templ)
		if err != nil {
			fmt.Printf("Input JSON file unmarshal error: %s\n", err.Error())
			os.Exit(1)
		}
	case "yaml":
		err = yaml.Unmarshal([]byte(inData), &templ)
		if err != nil {
			fmt.Printf("Input YAML file unmarshal error: %s\n", err.Error())
			os.Exit(1)
		}
	}

	var outData []byte

	switch sOut[len(sOut)-1] {
	case "xml":
		outData, err = xml.MarshalIndent(templ, "", "    ")
		if err != nil {
			fmt.Printf("Input XML file marshal error: %s\n", err.Error())
			os.Exit(1)
		}
		// pair replacement to match Zabbix XML style
		r := strings.NewReplacer("&#34;", "\"",
			"&#39;", "'",
			"&#xD;", "&#13;",
			"&#xA;", "\n",
			"<url></url>", "<url/>",
			"<params></params>", "<params/>",
			"<application></application>", "<application/>",
			"<items></items>", "",
                        "<macros></macros>", "",
			"<screens></screens>", "",
			"<mappings></mappings>", "",
			"<triggers></triggers>", "",
			"<conditions></conditions>", "",
			"<master_item></master_item>", "",
			"<dependencies></dependencies>", "",
                        "<applications></applications>", "",
			"<preprocessing></preprocessing>", "",
                        "<discovery_rules></discovery_rules>", "",
			"<item_prototypes></item_prototypes>", "",
			"<graph_prototypes></graph_prototypes>", "",
                        "<trigger_prototypes></trigger_prototypes>", "",
		)
		re := regexp.MustCompile("(?m)^\\s*$[\r\n]*")
		outData = []byte(xml.Header + re.ReplaceAllString(r.Replace(string(outData)), "") + "\n")
	case "json":
		outData, err = json.MarshalIndent(templ, "", "    ")
		if err != nil {
			fmt.Printf("Input JSON file marshal error: %s\n", err.Error())
			os.Exit(1)
		}
	case "yaml":
		outData, err = yaml.Marshal(templ)
		if err != nil {
			fmt.Printf("Input YAML file unmarshal error: %s\n", err.Error())
			os.Exit(1)
		}
	}

	f, err := os.Create(strings.Join(sOut, "."))
	if err != nil {
		fmt.Printf("Can't create output file, error: %s\n", err.Error())
		os.Exit(1)

	}
	defer f.Close()

	_, err = f.WriteString(string(outData))
	if err != nil {
		fmt.Printf("Can't write output file, error: %s\n", err.Error())
		os.Exit(1)
	}
	f.Sync()
}
