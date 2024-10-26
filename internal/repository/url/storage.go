package url

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/RomanAgaltsev/urlcut/internal/model"
)

type storage struct {
	f string
	m map[string]*model.URL
}

func newStorage(storagePath string) (*storage, error) {
	return &storage{
		f: storagePath,
	}, nil
}

func (s *storage) save() error {
	file, err := os.OpenFile(s.f, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	for _, v := range s.m {
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

func (s *storage) restore() error {
	file, err := os.OpenFile(s.f, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	s.m = make(map[string]*model.URL)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		var u model.URL
		if err := json.Unmarshal(data, &u); err != nil {
			return err
		}
		s.m[u.ID] = &u
	}

	return file.Close()
}
