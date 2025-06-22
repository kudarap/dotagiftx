package cache

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("cache set:", path)

	d, err := newData(val, expr)
	if err != nil {
		return err
	}
	fmt.Println(string(d))

	//if err = os.WriteFile(path, d, 0666); err != nil {
	//	return err
	//}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if _, err = f.Write(d); err != nil {
		return err
	}

	fmt.Println("WriteFile!")
	return nil
}

func Del(key string) error {
	path := filename(key)
	return os.Remove(path)
}

func newData(val interface{}, d time.Duration) ([]byte, error) {
	t := time.Now().Add(d).Unix()
	c := &data{val, t}
	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func filename(key string) string {
	return filepath.Join(os.TempDir(), dirname, key)
}
