package utils

import (
	"fmt"
	"os"
)

func WriteKeyValueToFile(filename, key, value string) (int64, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("error opening file: %v", err)
	}

	position, err := file.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, fmt.Errorf("error seeking file: %v", err)
	}

	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%s=%s\n", key, value))

	if err != nil {
		return 0, fmt.Errorf("error writing to file: %v", err)
	}
	return position, nil
}

func ReadFromFileAtPosition(filename string, position int64) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}

	defer file.Close()
	_, err = file.Seek(position, 0)
	if err != nil {
		return "", fmt.Errorf("error seeking file: %v", err)
	}

	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	return string(buf[:n]), nil
}

func MarkDelete(filename, key string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%s=%s\n", key, "*"))
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}
