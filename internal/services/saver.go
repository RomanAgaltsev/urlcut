package services

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/RomanAgaltsev/urlcut/internal/model"
)

type StateSaver struct {
	filename string
}

func NewStateSaver(filename string) *StateSaver {
	return &StateSaver{
		filename: filename,
	}
}

func (s *StateSaver) SaveState(state map[string]*model.URL) error {
	file, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	for _, v := range state {
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

func (s *StateSaver) RestoreState() (map[string]*model.URL, error) {
	file, err := os.OpenFile(s.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	state := make(map[string]*model.URL)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		var u model.URL
		if err := json.Unmarshal(data, &u); err != nil {
			return nil, err
		}
		state[u.ID] = &u
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return state, nil
}
