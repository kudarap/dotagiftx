package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const dirname = "gokitcache"

func init() {
	// This will create cache dir.
	_ = os.MkdirAll(filename(""), 0777)
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
	path := filename(key)
	b, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if b == nil {
		return nil, nil
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
	path := filename(key)
	err := os.WriteFile(path, newData(val, expr), 0666)
	if err != nil {
		return err
	}

	return nil
}

func Del(key string) error {
	path := filename(key)
	return os.Remove(path)
}

func filename(key string) string {
	return filepath.Join(os.TempDir(), dirname, key)
}
