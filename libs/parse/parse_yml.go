package parse

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ParseYmlFile parse a .yml file into a provided struct
func ParseYmlFile(file string, cfg interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	return nil
}

// ShouldParseYmlFile  parse a .yml file and panics if any error happens
func ShouldParseYmlFile(file string, cfg interface{}) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		panic(err)
	}
}
