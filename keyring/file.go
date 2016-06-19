package keyring

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func init() {
	supportedBackends[FileBackend] = opener(func(name string) (Keyring, error) {
		if name == "" {
			name = "file"
		}

		return &fileKeyring{name: name}, nil
	})
}

type fileKeyring struct {
	name string
}

func (k *fileKeyring) Get(key string) (Item, error) {

	items, err := Items()
	for _, item := range items {
		if item.Key == key {
			return item, nil
		}
	}
	return Item{}, err
}

func (k *fileKeyring) Set(item Item) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	// todo: read old entries, append, only then write
	err = ioutil.WriteFile("/tmp/dat", data, 0644)

	if err != nil {
		return err
	}

	return nil
}

func (k *fileKeyring) Remove(key string) error {
	return nil
}

func (k *fileKeyring) Keys() ([]string, error) {
	keys := []string{}
	if file, err := os.Open("/tmp/dat"); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			item := Item{}
			err := json.Unmarshal([]byte(line), &item)
			if err != nil {
				return keys, err
			}
			fmt.Println(item)
			keys = append(keys, item.Key)
		}
	} else {
		return []string{}, err
	}

	return keys, nil
}

func Items() ([]Item, error) {
	items := []Item{}
	if file, err := os.Open("/tmp/dat"); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			item := Item{}
			err := json.Unmarshal([]byte(line), &item)
			if err != nil {
				return items, err
			}
			items = append(items, item)
		}
	} else {
		return []Item{}, err
	}
	return items, nil
}
