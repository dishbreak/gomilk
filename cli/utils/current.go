package utils

import (
	"io/ioutil"
	"os/user"
	"path"
)

type currentRecord struct {
	Name     string
	Filepath string
}

/*
CurrentRecord represents the "current" value of something on the filesystem.

Gomilk uses these records to keep track of undoable operations.
*/
type CurrentRecord interface {
	Get() (string, error)
	Set(input string) error
}

func (c *currentRecord) Get() (string, error) {

	buf, err := ioutil.ReadFile(c.Filepath)
	if err != nil {
		return "", err
	}

	return string(buf), err
}

func (c *currentRecord) Set(input string) error {
	return ioutil.WriteFile(c.Filepath, []byte(input), 0600)
}

/*
NewCurrentRecord generates a CurrentRecord implementation that the caller can use to store data.
*/
func NewCurrentRecord(name string) (CurrentRecord, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	filePath := path.Join(currentUser.HomeDir, ".gomilk", "current", "current_"+name)

	c := &currentRecord{
		Name:     name,
		Filepath: filePath,
	}
	if err := mkdirIfNotExists(path.Dir(c.Filepath)); err != nil {
		return nil, err
	}

	return c, nil
}
