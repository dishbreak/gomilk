package utils

import (
	"fmt"
	"os"
	"os/user"
	"path"
)

/*
GetGomilkFileName will give you a filename for a file inside the .gomilk dir within the user's home directory.
It's your responsibility as the caller to create this file!
*/
func GetGomilkFile(filename string) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	gomilkDir := path.Join(user.HomeDir, ".gomilk")
	err = mkdirIfNotExists(gomilkDir)
	if err != nil {
		return "", err
	}

	return path.Join(gomilkDir, filename), nil

}

func mkdirIfNotExists(dirPath string) error {
	if stat, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	} else if mode := stat.Mode(); !mode.IsDir() {
		return fmt.Errorf("needed to create dir '%s' but it exists already as a file", dirPath)
	}
	return nil
}
