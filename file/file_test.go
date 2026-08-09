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

		want := generateHashName()
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

func TestLocal_resolvePath(t *testing.T) {
	local := Local{
		allowedTypes: []string{"image/jpeg"},
		sizeLimit:    kbSize * 1000,
		saveDir:      "/images",
	}

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "empty",
			path:    "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "simple filename",
			path:    "test.jpg",
			want:    "/images/test.jpg",
			wantErr: false,
		},
		{
			name:    "relative path traversal",
			path:    "../test.jpg",
			want:    "",
			wantErr: true,
		},
		{
			name:    "absolute path",
			path:    "/tmp/test.jpg",
			want:    "",
			wantErr: true,
		},
		{
			name:    "multiple directory traversal",
			path:    "../../../test.jpg",
			want:    "",
			wantErr: true,
		},
		{
			name:    "path with dots",
			path:    "./test.jpg",
			want:    "/images/test.jpg",
			wantErr: false,
		},
		{
			name:    "path with multiple dots",
			path:    ".../test.jpg",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty path component",
			path:    "test//file.jpg",
			want:    "/images/test/file.jpg",
			wantErr: false,
		},
		{
			name:    "path with spaces",
			path:    "test file.jpg",
			want:    "/images/test file.jpg",
			wantErr: false,
		},
		{
			name:    "path with special characters",
			path:    "test@#$%^&*().jpg",
			want:    "",
			wantErr: true,
		},
		{
			name:    "path with unicode characters",
			path:    "test_测试.jpg",
			want:    "",
			wantErr: true,
		},
		{
			name:    "path with encoded characters",
			path:    "test%20file.jpg",
			want:    "",
			wantErr: true,
		},
		{
			name:    "path with backslashes",
			path:    "test\\file.jpg",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := local.resolvePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("cleanPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cleanPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func errCheck(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
