package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

/*
Identifiable represents any object that has a canonical id and a friendly name.
The id must be unique, the name need not be.
*/
type Identifiable interface {
	ID() string
	Name() string
}

/*
Cache allows a gomilk command the ability to store a list of API objects.
This allows the user to reference API objects in a subsequent command.
*/
type Cache interface {
	Filename() string
}

type cache struct {
	name     string
	filename string
}

type identifierRecord struct {
	Index   int
	RawID   string
	RawName string
}

type cacheContents []identifierRecord

func (i identifierRecord) ID() string {
	return i.RawID
}

func (i identifierRecord) Name() string {
	return i.RawName
}

func (i identifierRecord) String() string {
	return fmt.Sprintf("<idx: %d id: %s name: %s>", i.Index, i.ID(), i.Name())
}

/*
NewCache creates a new cache on the filesystem with the given identifier.
This identifier must be unique!
*/
func NewCache(identifier string) (Cache, error) {
	f, err := GetGomilkFile(identifier + "_cache")
	c := &cache{identifier, f}

	return c, err
}

func (c *cache) Filename() string {
	return c.filename
}

func (c *cache) Update(items []interface{}) error {

	data, err := json.Marshal(items)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.filename, data, 0600)

	return err
}

func (c *cache) Get(items []interface{}) ([]interface{}, error) {
	data, err := ioutil.ReadFile(c.filename)
	if err != nil {
		return items, err
	}

	err = json.Unmarshal(data, &items)
	if err != nil {
		return items, err
	}

	return items, err
}
