package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"
)

func main() {
	configBytes, err := ioutil.ReadAll(os.Stdin)
	CheckErr(err, "Error reading stdin:", err)

	source := GenerateSource(configBytes)
	fmt.Print(source)
}

func GenerateSource(configBytes []byte) string {
	var config map[string]interface{}
	err := yaml.Unmarshal(configBytes, &config)
	CheckErr(err, "Error parsing alias config", err)

	aliases := buildAST(config)
	return PrettyPrint(aliases)
}

func buildAST(config map[string]interface{}) []*Alias {
	subAliases := []*Alias{}

	for k, v := range config {
		subAlias := &Alias{Identifier: k}

		t := reflect.TypeOf(v).Kind()
		switch t {
		case reflect.String:
			command, ok := v.(string)
			CheckOk(ok, "Unable to convert to string:", v)
			subAlias.Command = command

		case reflect.Map:
			subAliasConfig, ok := v.(map[string]interface{})
			CheckOk(ok, "Unable to convert to map:", v)

			if commandObject, ok := subAliasConfig["$"]; ok {
				command, ok := commandObject.(string)
				CheckOk(ok, "Unable to convert to string:", commandObject)
				subAlias.Command = command
				delete(subAliasConfig, "$")
			}
			subAlias.SubAliases = buildAST(subAliasConfig)

		default:
			fmt.Println("Alias config must be string or map:", v)
			os.Exit(1)
		}
		subAliases = append(subAliases, subAlias)
	}

	return subAliases
}
