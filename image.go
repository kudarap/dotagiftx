package dotagiftx

import (
	"context"
	"io"

	"github.com/kudarap/dotagiftx/file"
)

// Image error types.
const (
	ImageErrNotFound Errors = iota + itemErrorIndex
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
		// Save saves file and returns a file name.
		Save(r io.Reader) (filename string, err error)

		// SaveWithName saves file with pre-defined base name.
		SaveWithName(r io.Reader, baseName string) (filename string, err error)

		// Get return file path base on file name.
		Get(filename string) (path string, err error)

		// Delete uploaded file base on file name.
		Delete(filename string) error

		// Dir returns save path location.
		Dir() string
	}
)

// NewImageService returns a new Image service.
func NewImageService(fm FileManager) ImageService {
	return &imageService{fm}
}

type imageService struct {
	fileMgr FileManager
}

func (s *imageService) Upload(ctx context.Context, r io.Reader) (fileID string, err error) {
	if au := AuthFromContext(ctx); au == nil {
		err = AuthErrNoAccess
		return
	}

	fileID, err = s.fileMgr.Save(r)
	if err != nil {
		err = NewXError(ImageErrUpload, err)
		return
	}

	return fileID, nil
}

func (s *imageService) Thumbnail(fileID string, width, height uint) (path string, err error) {
	f, err := s.Image(fileID)
	if err != nil {
		return
	}

	t, err := file.Thumbnail(f, width, height)
	if err != nil {
		err = NewXError(ImageErrThumbnail, err)
		return
	}

	return t, nil
}

func (s *imageService) Image(fileID string) (path string, err error) {
	path, err = s.fileMgr.Get(fileID)
	if err != nil {
		err = NewXError(ImageErrNotFound, err)
		return
	}

	return path, nil
}

func (s *imageService) Delete(ctx context.Context, fileID string) error {
	if au := AuthFromContext(ctx); au == nil {
		return AuthErrNoAccess
	}

	return s.fileMgr.Delete(fileID)
}
