package embedded

import (
	"embed"
	"io/fs"
)

//go:embed md/*
var content embed.FS

// ReadFile reads the contents of an embedded file and returns it as a string.
func ReadFile(filename string) (string, error) {
	file, err := content.Open("md/" + filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := fs.ReadFile(content, "md/"+filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
