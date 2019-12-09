package parse

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// YmlFile parse a .yml file into a provided struct
func YmlFile(file string, cfg interface{}) error {
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

// ShouldYmlFile  parse a .yml file and panics if any error happens
func ShouldYmlFile(file string, cfg interface{}) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		panic(err)
	}
}
