package buffer

import "os"

func createFileIfNotExists(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if file == nil || err != nil {
		return nil, err
	}

	return file, nil
}
