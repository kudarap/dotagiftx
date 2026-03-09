package file

import (
	"os"
	"strings"
	"testing"
)

func TestLocal_Save(t *testing.T) {
	local := Local{
		allowedTypes: []string{"image/jpeg"},
		sizeLimit:    kbSize * 1000,
		saveDir:      os.TempDir(),
	}

	t.Run("file base name", func(t *testing.T) {
		image, err := os.Open("testdata/profile.jpg")
		errCheck(t, err)

		want := generateSha1Name()
		name, err := local.SaveWithName(image, want)
		errCheck(t, err)
		if !strings.HasPrefix(name, want) {
			t.Errorf("SaveWithName should want prefix %s, got: %s", want, name)
		}
	})
	t.Run("file extension", func(t *testing.T) {
		image, err := os.Open("testdata/profile.jpg")
		errCheck(t, err)

		want := ".jpg"
		name, err := local.Save(image)
		errCheck(t, err)
		if !strings.HasSuffix(name, ".jpg") {
			t.Errorf("Save() should want suffix %s, got: %s", want, name)
		}
	})
}

func errCheck(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
