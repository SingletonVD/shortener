package repository

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/SingletonVD/shortener/internal/model"
)

type DiskLinkRepository struct {
	fileStoragePath string
	links           map[string]string
	lock            sync.Mutex
}

func NewDiskLinkRepository(fileStoragePath string) (*DiskLinkRepository, error) {
	repo := DiskLinkRepository{
		fileStoragePath: fileStoragePath,
		links:           make(map[string]string),
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

	storage.links[link.Short] = link.FullLink
	err := storage.dumpAll()

	if err != nil {
		delete(storage.links, link.Short)
		return false, err
	}

	return true, nil
}

func (storage *DiskLinkRepository) FindLink(shortLink string) (*model.ShortenedLink, error) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	fullLink, found := storage.links[shortLink]

	if !found {
		return nil, nil
	}

	return &model.ShortenedLink{
		Short:    shortLink,
		FullLink: fullLink,
	}, nil
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

	var persistentLinks []model.PersistedShortenedLink

	uuid := 1
	for short, long := range storage.links {
		persistentLink := model.PersistedShortenedLink{
			UUID: strconv.Itoa(uuid),
			ShortenedLink: model.ShortenedLink{
				Short:    short,
				FullLink: long,
			},
		}
		persistentLinks = append(persistentLinks, persistentLink)
		uuid++
	}

	return encoder.Encode(persistentLinks)
}

func (storage *DiskLinkRepository) restoreState() error {
	file, err := os.OpenFile(storage.fileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	var links []model.ShortenedLink
	err = decoder.Decode(&links)

	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	for _, link := range links {
		storage.links[link.Short] = link.FullLink
	}

	return nil
}
