package repository

import (
	"encoding/json"
	"errors"
	"io"
	"maps"
	"os"
	"slices"
	"sync"

	"github.com/SingletonVD/shortener/internal/model"
)

type DiskLinkRepository struct {
	fileStoragePath string
	links           map[string]model.PersistedShortenedLink
	lock            sync.Mutex
	lastId          int // что-то типа автоинкремента
}

func NewDiskLinkRepository(fileStoragePath string) (*DiskLinkRepository, error) {
	repo := DiskLinkRepository{
		fileStoragePath: fileStoragePath,
		links:           make(map[string]model.PersistedShortenedLink),
	}
	err := repo.restoreState()

	if err != nil {
		return nil, err
	}

	return &repo, err
}

func (storage *DiskLinkRepository) SaveIfAvailable(link model.ShortenedLink) (bool, error) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	_, found := storage.links[link.Short]
	if found {
		return false, nil
	}

	persistedLink := model.PersistedShortenedLink{
		ShortenedLink: link,
		UUID:          storage.lastId + 1,
	}
	storage.links[link.Short] = persistedLink
	err := storage.dumpAll()

	if err != nil {
		delete(storage.links, link.Short)
		return false, err
	}

	storage.lastId = persistedLink.UUID
	return true, nil
}

func (storage *DiskLinkRepository) FindLink(shortLink string) (*model.ShortenedLink, error) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	persistedLink, found := storage.links[shortLink]

	if !found {
		return nil, nil
	}

	return &persistedLink.ShortenedLink, nil
}

// в примере файл хранилища - JSON массив, поэтому делаю dump всего,
// чтобы не усложнять с вставкой именно в конец с сохранением структуры,
// но вообще я бы лучше сделал JSON Lines с вставкой в конец файла
func (storage *DiskLinkRepository) dumpAll() error {
	file, err := os.OpenFile(storage.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	persistentLinks := slices.Collect(maps.Values(storage.links))
	return encoder.Encode(persistentLinks)
}

func (storage *DiskLinkRepository) restoreState() error {
	file, err := os.OpenFile(storage.fileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	var links []model.PersistedShortenedLink
	err = decoder.Decode(&links)

	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	for _, link := range links {
		storage.links[link.Short] = link
		storage.lastId = max(storage.lastId, link.UUID)
	}

	return nil
}
