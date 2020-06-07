package file

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const kbSize = 1000

type Local struct {
	saveDir      string
	sizeLimit    int64
	allowedTypes []string
}

// New creates a new instance of Local file manager.
func New(saveDir string, sizeLimit int, allowedTypes []string) *Local {
	sizeLimit *= kbSize
	return &Local{saveDir, int64(sizeLimit), allowedTypes}
}

// Save saves bytes into file and returns an unique filename.
func (l *Local) Save(r io.Reader) (name string, err error) {
	// Check empty save saveDir.
	if strings.TrimSpace(l.saveDir) == "" {
		err = errors.New("file save saveDir required")
		return
	}
	// Make upload path writable.
	dir := strings.TrimSuffix(l.Dir(), "/") + "/"
	_ = os.Mkdir(dir, os.ModePerm)

	// Get bytes content from Reader.
	buf := new(bytes.Buffer)
	size, err := buf.ReadFrom(r)
	if err != nil {
		return
	}
	data := buf.Bytes()

	// Validate file size.
	if size >= l.sizeLimit {
		err = fmt.Errorf("file size limit reached %d of %d KB", size/kbSize, l.sizeLimit/kbSize)
		return
	}
	// Validate file type.
	cType, err := l.getType(data)
	if err != nil {
		return
	}

	// Generate name from file check sum.
	_, name, err = checksum(buf)
	if err != nil {
		return
	}
	// Compose file name with extension.
	ext, err := mime.ExtensionsByType(cType)
	if err != nil {
		return
	}
	name += normalizeExt(ext[0])

	// Create file inside save saveDir.
	dst := filepath.Join(l.Dir(), name)
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer out.Close()

	// Write contents to file.
	_, err = out.Write(data)
	if err != nil {
		return
	}

	err = nil
	return
}

// GetFile get file path base on file name and its existence.
func (l *Local) Get(name string) (path string, err error) {
	path = filepath.Join(l.Dir(), name)
	// Check actual file existence.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", err
	}

	return
}

// Delete uploaded file base on file name.
func (l *Local) Delete(name string) error {
	p := filepath.Join(l.Dir(), name)
	if err := os.Remove(p); err != nil {
		return err
	}

	return nil
}

// Dir returns save path location.
func (l *Local) Dir() string {
	return l.saveDir
}

func (l *Local) getType(data []byte) (string, error) {
	t := http.DetectContentType(data)
	for _, tt := range l.allowedTypes {
		if strings.HasPrefix(t, tt) {
			return t, nil
		}
	}

	// No hit return error.
	return "", fmt.Errorf("file type '%s' not allowed in %s", t, l.allowedTypes)
}

func checksum(r io.Reader) (io.Reader, string, error) {
	var b bytes.Buffer

	h := sha1.New()
	if _, err := io.Copy(&b, io.TeeReader(r, h)); err != nil {
		return &b, "", err
	}

	s := fmt.Sprintf("%d", time.Now().Nanosecond())
	h.Write([]byte(s))
	sum := h.Sum(nil)
	return &b, hex.EncodeToString(sum), nil
}

func normalizeExt(ext string) string {
	switch ext {
	case ".jpeg":
		return ".jpg"
	}

	return ext
}
