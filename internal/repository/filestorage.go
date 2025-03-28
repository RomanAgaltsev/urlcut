package repository

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/RomanAgaltsev/urlcut/internal/model"
)

// readFromFile считывает данные с URL из файла в мапу и возвращает её.
func readFromFile(fileStoragePath string) (map[string]*model.URL, error) {
	// Открываем файл на чтение
	file, err := os.OpenFile(fileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	// Создаем мапу
	m := make(map[string]*model.URL)

	// Читаем данные файла сканером
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		var u model.URL
		if err := json.Unmarshal(data, &u); err != nil {
			return nil, err
		}
		m[u.ID] = &u
	}

	// Возвращаем мапу
	return m, nil
}

// writeToFile сохраняет переданную мапу с URL в файл по переданному пути.
func writeToFile(fileStoragePath string, m map[string]*model.URL) error {
	// Открываем файл на запись
	file, err := os.OpenFile(fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	// Писать в файл будем через буфер
	writer := bufio.NewWriter(file)
	for _, v := range m {
		data, errj := json.Marshal(*v)
		if errj != nil {
			return errj
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

	// Закрываем файл
	return file.Close()
}
