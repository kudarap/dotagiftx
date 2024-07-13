package dgx

import (
	"context"
	"io"
)

// Image error types.
const (
	ImageErrNotFound Errors = iota + 3000
	ImageErrUpload
	ImageErrThumbnail
)

// sets error text definition.
func init() {
	appErrorText[ImageErrNotFound] = "image file not found"
	appErrorText[ImageErrUpload] = "image file upload error"
	appErrorText[ImageErrThumbnail] = "thumbnail image processing error"
}

type (
	// Image represents file image information.
	Image struct {
		FileID  string `json:"file_id"    db:"file_id,omitempty" valid:"required"`
		Caption string `json:"caption"    db:"caption,omitempty"`
	}

	// ImageService provides access image services.
	ImageService interface {
		// Upload saves image details and actual file to local file system.
		Upload(context.Context, io.Reader) (fileID string, err error)

		// Image returns image details by id.
		Image(fileID string) (path string, err error)

		// Thumbnail downscales an image preserving its aspect ratio to the maximum dimensions.
		// It will return the original image if original sizes are smaller than the provided dimensions.
		Thumbnail(fileID string, width, height uint) (path string, err error)

		// Delete purges image record and from local file system.
		Delete(ctx context.Context, fileID string) error
	}

	// FileManager defines operation for file on local file system.
	FileManager interface {
		// Upload save file and returns a file name.
		Save(r io.Reader) (filename string, err error)

		// Upload save file with pre-defined base name.
		SaveWithName(r io.Reader, baseName string) (filename string, err error)

		// Get get file path base on file name.
		Get(filename string) (path string, err error)

		// Delete uploaded file base on file name.
		Delete(filename string) error

		// Dir returns save path location.
		Dir() string
	}
)
