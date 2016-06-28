package keyring

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

func init() {
	supportedBackends[FileBackend] = opener(func(name string) (Keyring, error) {
		return &fileKeyring{}, nil
	})
}

type fileKeyring struct {
	Dir string
}

func (k *fileKeyring) dir() (string, error) {
	dir := k.Dir
	if dir == "" {
		usr, err := user.Current()
		if err != nil {
			return dir, err
		}
		dir = usr.HomeDir + "/.awsvault/keys/"
	}

	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.MkdirAll(dir, 0700)
	} else if err != nil && !stat.IsDir() {
		err = fmt.Errorf("%s is a file, not a directory", dir)
	}

	return dir, nil
}

func (k *fileKeyring) Get(key string) (Item, error) {
	dir, err := k.dir()
	if err != nil {
		return Item{}, err
	}

	bytes, err := ioutil.ReadFile(filepath.Join(dir, key))
	if os.IsNotExist(err) {
		return Item{}, ErrKeyNotFound
	} else if err != nil {
		return Item{}, err
	}
	item := Item{}
	err = json.Unmarshal(bytes, &item)

	return item, err
}

func (k *fileKeyring) Set(i Item) error {
	bytes, err := json.Marshal(i)
	if err != nil {
		return err
	}

	dir, err := k.dir()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(dir, i.Key), bytes, 0600)
}

func (k *fileKeyring) Remove(key string) error {
	dir, err := k.dir()
	if err != nil {
		return err
	}

	return os.Remove(filepath.Join(dir, key))
}

func (k *fileKeyring) Keys() ([]string, error) {
	dir, err := k.dir()
	if err != nil {
		return nil, err
	}

	var keys = []string{}
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		keys = append(keys, f.Name())
	}

	return keys, nil
}
