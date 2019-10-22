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

	if templVersion.Version != "4.4" && templVersion.Version != "4.0" {
		fmt.Printf("Unsupported Zabbix template version %s, only versions 4.4/4.0 are supported\n", templVersion.Version)
		os.Exit(1)
	}

	var templ interface{}
	if templVersion.Version == "4.4" {
		//v4.4
		templ = &v44.ZabbixExport{}
	} else {
		// v4.0
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
		if templVersion.Version == "4.4" {
			// pair replacement to match Zabbix XML style
			r := strings.NewReplacer("&#34;", "\"",
				"&#39;", "'",
				"&#xD;", "&#13;",
				"&#xA;", "\n",
				"<url></url>", "<url/>",
				"<params></params>", "<params/>",
				"<application></application>", "<application/>",
				"<items></items>", "",
				"<groups></groups>", "",
				"<graphs></graphs>", "",
				"<macros></macros>", "",
				"<screens></screens>", "",
				"<mappings></mappings>", "",
				"<triggers></triggers>", "",
				"<templates></templates>", "",
				"<conditions></conditions>", "",
				"<value_maps></value_maps>", "",
				"<master_item></master_item>", "",
				"<dependencies></dependencies>", "",
				"<applications></applications>", "",
				"<preprocessing></preprocessing>", "",
				"<lld_macro_paths></lld_macro_paths>", "",
				"<discovery_rules></discovery_rules>", "",
				"<item_prototypes></item_prototypes>", "",
				"<graph_prototypes></graph_prototypes>", "",
				"<trigger_prototypes></trigger_prototypes>", "",
				"<application_prototypes></application_prototypes>", "",
			)
			re := regexp.MustCompile(`(?m)^\s*$[\r\n]*`)
			re2 := regexp.MustCompile(`(?m)^[\s]*<valuemap>[\s]*</valuemap>[\s]*$`)
			re3 := regexp.MustCompile(`(?m)^[\s]*<filter>[\s]*</filter>[\s]*$`)
			outData = []byte(xml.Header + re.ReplaceAllString(re2.ReplaceAllString(re3.ReplaceAllString(r.Replace(string(outData)), ""), ""), ""))
		}
		if templVersion.Version == "4.0" {
			// pair replacement to match Zabbix XML style
			r := strings.NewReplacer("&#34;", "\"",
				"&#39;", "'",
				"&#xD;", "&#13;",
				"&#xA;", "\n",
				"<key></key>", "",
				"<name></name>", "",
				"<error_handler></error_handler>", "",
				"<error_handler_params></error_handler_params>", "",
				"<url></url>", "<url/>",
				"<params></params>", "<params/>",
				"<application></application>", "<application/>",
				"<items></items>", "<items></items>",
				"<groups></groups>", "<groups/>",
				"<graphs></graphs>", "",
				"<macros></macros>", "<macros/>",
				"<screens></screens>", "<screens/>",
				"<mappings></mappings>", "",
				"<triggers></triggers>", "",
				"<templates></templates>", "<templates/>",
				"<conditions></conditions>", "<conditions/>",
				"<logtimefmt></logtimefmt>", "<logtimefmt/>",
				"<description></description>", "<description/>",
				"<snmp_community></snmp_community>", "<snmp_community/>",
				"<snmp_oid></snmp_oid>", "<snmp_oid/>",
				"<value_maps></value_maps>", "<value_maps/>",
				"<master_item></master_item>", "<master_item/>",
				"<dependencies></dependencies>", "<dependencies/>",
				"<applications></applications>", "<applications/>",
				"<preprocessing></preprocessing>", "<preprocessing/>",
				"<lld_macro_paths></lld_macro_paths>", "<lld_macro_paths/>",
				"<discovery_rules></discovery_rules>", "<discovery_rules/>",
				"<item_prototypes></item_prototypes>", "<item_prototypes/>",
				"<graph_prototypes></graph_prototypes>", "<graph_prototypes/>",
				"<trigger_prototypes></trigger_prototypes>", "",
				"<application_prototypes></application_prototypes>", "<application_prototypes/>",
				"<recovery_expression></recovery_expression>", "<recovery_expression/>",
				"<allowed_hosts></allowed_hosts>", "<allowed_hosts/>",
				"<units></units>", "<units/>",
				"<username></username>", "<username/>",
				"<password></password>", "<password/>",
				"<snmpv3_contextname></snmpv3_contextname>", "<snmpv3_contextname/>",
				"<snmpv3_securityname></snmpv3_securityname>", "<snmpv3_securityname/>",
				"<snmpv3_authpassphrase></snmpv3_authpassphrase>", "<snmpv3_authpassphrase/>",
				"<snmpv3_privpassphrase></snmpv3_privpassphrase>", "<snmpv3_privpassphrase/>",
				"<publickey></publickey>", "<publickey/>",
				"<privatekey></privatekey>", "<privatekey/>",
				"<port></port>", "<port/>",
				"<correlation_tag></correlation_tag>", "<correlation_tag/>",
				"<ssl_cert_file></ssl_cert_file>", "<ssl_cert_file/>",
				"<ssl_key_file></ssl_key_file>", "<ssl_key_file/>",
				"<ssl_key_password></ssl_key_password>", "<ssl_key_password/>",
				"<ipmi_sensor></ipmi_sensor>", "<ipmi_sensor/>",
				"<jmx_endpoint></jmx_endpoint>", "<jmx_endpoint/>",
			)
			re := regexp.MustCompile(`(?m)^\s*$[\r\n]*`)
			re2 := regexp.MustCompile(`(?m)^<valuemap>[\s]*</valuemap>[\s]*$`)
			re3 := regexp.MustCompile(`(?m)^[\s]*<filter>[\s]*</filter>[\s]*$`)
			outData = []byte(xml.Header + re.ReplaceAllString(re2.ReplaceAllString(re3.ReplaceAllString(r.Replace(string(outData)), ""), "<valuemap/>"), ""))
		}
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
	err = f.Sync()
	if err != nil {
		fmt.Printf("Can't sync output file, error: %s\n", err.Error())
		os.Exit(1)
	}
}
