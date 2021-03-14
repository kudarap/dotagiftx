package cache

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var cacheDir string

func init() {
	cacheDir = os.TempDir()
}

type data struct {
	Payload interface{}
	Expr    int64
}

func (d *data) isExpired() bool {
	return time.Now().Unix() > d.Expr
}

func newData(val interface{}, d time.Duration) []byte {
	t := time.Now().Add(d).Unix()
	c := &data{val, t}
	b, _ := json.Marshal(c)
	return b
}

func Get(key string) (val interface{}, err error) {
	path := filepath.Join(cacheDir, key)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	d := new(data)
	if err = json.Unmarshal(b, d); err != nil {
		return nil, err
	}

	if d.isExpired() {
		return nil, Del(key)
	}

	return d.Payload, nil
}

func Set(key string, val interface{}, expr time.Duration) error {
	path := filepath.Join(cacheDir, key)
	err := ioutil.WriteFile(path, newData(val, expr), 0666)
	if err != nil {
		return err
	}

	log.Println("cache write", path)

	return nil
}

func Del(key string) error {
	path := filepath.Join(cacheDir, key)
	return os.Remove(path)
}
