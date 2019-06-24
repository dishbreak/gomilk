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
	// Update replaces the contents of the cache with the array of Identifiable objects.
	Update([]Identifiable) error
	// Get returns the current contents of the cache as an array of Identifiable objects.
	Get() ([]Identifiable, error)
}

type cache struct {
	name     string
	filename string
}

type identifierRecord struct {
	Index   int
	RawID   string
	RawName string
	Data    interface{}
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

func (c *cache) Update(items []Identifiable) error {
	output := make(cacheContents, len(items))
	for idx, item := range items {
		output[idx] = identifierRecord{idx, item.ID(), item.Name()}
	}

	data, err := json.Marshal(output)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.filename, data, 0600)

	return err
}

func (c *cache) Get() ([]Identifiable, error) {
	data, err := ioutil.ReadFile(c.filename)
	if err != nil {
		return []Identifiable{}, err
	}

	var contents []identifierRecord
	err = json.Unmarshal(data, &contents)
	if err != nil {
		return []Identifiable{}, err
	}

	result := make([]Identifiable, len(contents))
	for idx, elem := range contents {
		result[idx] = elem
	}

	return result, err
}
