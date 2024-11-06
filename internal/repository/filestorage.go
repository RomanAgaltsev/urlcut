package repository

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/RomanAgaltsev/urlcut/internal/model"
)

func readFromFile(fileStoragePath string) (map[string]*model.URL, error) {
	file, err := os.OpenFile(fileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	m := make(map[string]*model.URL)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		var u model.URL
		if err := json.Unmarshal(data, &u); err != nil {
			return nil, err
		}
		m[u.ID] = &u
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return m, nil
}

func writeToFile(fileStoragePath string, m map[string]*model.URL) error {
	file, err := os.OpenFile(fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	for _, v := range m {
		data, err := json.Marshal(*v)
		if err != nil {
			return err
		}
		if _, err = writer.Write(data); err != nil {
			return err
		}
		if err = writer.WriteByte('\n'); err != nil {
			return err
		}
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	return file.Close()
}
