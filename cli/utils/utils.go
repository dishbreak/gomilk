package utils

import (
	"fmt"
	"os"
	"os/user"
	"path"
)

/*
GetGomilkFileName will give you a pointer to a file inside the .gomilk dir within the user's home directory.

Be warned! This will erase any existing file with the same name.
*/
func GetGomilkFile(filename string) (*os.File, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	gomilkDir := path.Join(user.HomeDir, ".gomilk")
	if stat, err := os.Stat(gomilkDir); os.IsNotExist(err) {
		err = os.Mkdir(gomilkDir, 0755)
		if err != nil {
			return nil, err
		}
	} else if mode := stat.Mode(); !mode.IsDir() {
		return nil, fmt.Errorf("needed to create dir '%s' but it exists already as a file", gomilkDir)
	}

	fileHandle, err := os.Create(path.Join(gomilkDir, filename))
	if err != nil {
		return nil, err
	}

	err = fileHandle.Chmod(0600)
	if err != nil {
		return nil, err
	}

	return fileHandle, nil
}
