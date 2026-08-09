package file

import (
	"bytes"
	"crypto/sha1" //nolint:gosec // legacy file name hashing for consistency, not security critical
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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

// SaveWithName saves bytes into file with pre-defined name.
func (l *Local) SaveWithName(r io.Reader, baseName string) (name string, err error) {
	return l.baseSave(r, baseName)
}

// Save saves bytes into file and returns a unique filename.
func (l *Local) Save(r io.Reader) (name string, err error) {
	return l.baseSave(r, generateHashName())
}

// Save saves bytes into file and returns an unique filename.
func (l *Local) baseSave(r io.Reader, baseName string) (name string, err error) {
	// Check empty save saveDir.
	if strings.TrimSpace(l.saveDir) == "" {
		err = errors.New("file save saveDir required")
		return
	}
	// Make upload path writable.
	dir := strings.TrimSuffix(l.Dir(), "/") + "/"
	_ = os.Mkdir(dir, 0o750)

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

	// Compose file name with extension.
	ext, err := mime.ExtensionsByType(cType)
	if err != nil {
		return
	}
	name = baseName + normalizeExt(ext[0])

	// Create file inside save saveDir.
	dst, err := l.resolvePath(name)
	if err != nil {
		return
	}
	out, err := os.Create(dst) //nolint:gosec // dst is validated by resolvePath against save dir
	if err != nil {
		return
	}
	defer func() { _ = out.Close() }()
	defer func() {
		if err := out.Close(); err != nil {
			slog.Error("closing file", "error", err)
		}
	}()

	// Write contents to file.
	_, err = out.Write(data)
	if err != nil {
		return
	}

	err = nil
	return
}

// Get gets file path base on file name and its existence.
func (l *Local) Get(name string) (path string, err error) {
	path, err = l.resolvePath(name)
	if err != nil {
		return "", err
	}

	// Check the actual file existence.
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return "", err
	}

	return
}

// Delete uploaded file base on file name.
func (l *Local) Delete(name string) error {
	p, err := l.resolvePath(name)
	if err != nil {
		return err
	}

	if err = os.Remove(p); err != nil {
		return err
	}

	return nil
}

// Dir returns save path location.
func (l *Local) Dir() string {
	return filepath.Clean(l.saveDir)
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

var validNamePattern = regexp.MustCompile(`^[A-Za-z0-9 ._/-]+$`)

func (l *Local) resolvePath(name string) (string, error) {
	// validate raw file name
	if filepath.IsAbs(name) || strings.HasPrefix(name, "..") || !validNamePattern.MatchString(name) {
		return "", fmt.Errorf("invalid file name '%s'", name)
	}

	// clean file name and path
	dir := l.Dir()
	path := filepath.Clean(filepath.Join(dir, name))

	// Check if the cleaned path is within the allowed directory
	// This prevents directory traversal attacks
	if !strings.HasPrefix(path, dir) {
		return "", fmt.Errorf("path %s is outside of allowed directory %s", name, l.Dir())
	}

	return path, nil
}

func generateHashName() string {
	h := sha1.New() //nolint:gosec // legacy file name hashing for consistency
	s := fmt.Sprintf("%d", time.Now().Nanosecond())
	h.Write([]byte(s))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func normalizeExt(ext string) string {
	switch ext {
	case ".jpe":
		fallthrough
	case ".jfif":
		fallthrough
	case ".jpeg":
		return ".jpg"
	}

	return ext
}
