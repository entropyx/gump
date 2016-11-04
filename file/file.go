package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type M map[interface{}]interface{}

type File struct {
	Attributes       M
	VersionLocations []string
	path             string
	extension        string
}

func Read(filename string) (*File, error) {
	file := &File{}
	file.setExtension(filename)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = unmarshal(bytes, file)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *File) Modify(version string, locations []string) error {
	m, err := f.generateNewContent(version, locations)
	if err != nil {
		return err
	}
	err = f.write(m)
	if err != nil {
		return err
	}
	f.Attributes = m
	return nil
}

func (f *File) generateNewContent(version string, locations []string) (M, error) {
	newContent := f.Attributes
	for _, location := range locations {
		var inf interface{}
		var prevInf interface{} //I'm using this var because couldn't set version directly
		positions := strings.Split(location, ".")
		for i, position := range positions {
			if inf == nil {
				inf = newContent[position]
				if inf == nil {
					err := fmt.Errorf("invalid field %s in %s", position, location)
					return nil, err
				}
			} else {
				prevInf = inf
				inf = inf.(map[interface{}]interface{})[position]
			}
			t := reflect.TypeOf(inf)
			switch t.String() {
			case "string", "int":
				if i != len(positions)-1 {
					err := fmt.Errorf("expected more fields in %s, but %s was the last one", location, position)
					return nil, err
				}
				prevInf.(map[interface{}]interface{})[position] = version
			case "map[string]interface {}", "map[interface {}]interface {}":

			default:
				return nil, fmt.Errorf("invalid data type")
			}
		}
	}
	return newContent, nil
}

func (f *File) write(m M) error {
	bytes, err := marshal(m, f.extension)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(f.path, bytes, 644)
	if err != nil {
		return err
	}
	return nil
}

func marshal(in interface{}, extension string) ([]byte, error) {
	var err error
	var bytes []byte
	switch extension {
	case "json":
		bytes, err = json.Marshal(in)
	case "yml", "yaml":
		bytes, err = yaml.Marshal(in)
	}
	return bytes, err
}

func (f *File) setExtension(filename string) {
	split := strings.Split(filename, ".")
	if len(split) > 1 {
		f.extension = split[len(split)-1]
	}
	f.path = filename
}

func unmarshal(in []byte, file *File) error {
	var fileContent map[interface{}]interface{}
	var err error
	switch file.extension {
	case "json":
		err = json.Unmarshal(in, &fileContent)
	case "yml", "yaml":
		err = yaml.Unmarshal(in, &fileContent)
	}
	if err != nil {
		return err
	}
	file.Attributes = fileContent
	return nil
}
