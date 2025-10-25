package utils

import "os"

func ReadFile(filePath string) ([]byte, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileData, nil
}

func WriteFile(filePath string, data []byte) error {
	err := os.WriteFile(filePath, data, 0755)
	if err != nil {
		return err
	}
	return nil
}
