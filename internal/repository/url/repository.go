package url

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/RomanAgaltsev/urlcut/internal/model"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
)

var _ repository.URLRepository = (*InMemoryRepository)(nil)

var ErrIDNotFound = fmt.Errorf("URL ID was not found in repository")

type InMemoryRepository struct {
	m map[string]*model.URL
	f string
	sync.RWMutex
}

func New(filename string) *InMemoryRepository {
	return &InMemoryRepository{
		m: make(map[string]*model.URL),
		f: filename,
	}
}

func (r *InMemoryRepository) Store(url *model.URL) error {
	r.Lock()
	defer r.Unlock()

	r.m[url.ID] = url

	return nil
}

func (r *InMemoryRepository) Get(id string) (*model.URL, error) {
	r.Lock()
	defer r.Unlock()

	if url, ok := r.m[id]; ok {
		return url, nil
	} else {
		return url, ErrIDNotFound
	}
}

func (r *InMemoryRepository) SaveState() error {
	file, err := os.OpenFile(r.f, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	for _, v := range r.m {
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

func (r *InMemoryRepository) RestoreState() error {
	file, err := os.OpenFile(r.f, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		var u model.URL
		if err := json.Unmarshal(data, &u); err != nil {
			return err
		}
		r.m[u.ID] = &u
	}

	return file.Close()
}
