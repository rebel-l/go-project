package workingdir

import "os"

var dir string

func Get() string {
	return dir
}

func Init() error {
	if err := detect(); err != nil {
		return err // TODO: ask for the folder
	}

	return nil
}

func detect() error {
	var err error
	dir, err = os.Getwd()
	return err
}
