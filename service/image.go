package service

import (
	"context"
	"io"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
	"github.com/kudarap/dotagiftx/gokit/file/image"
)

// NewAuth returns a new Image service.
func NewImage(fm dgx.FileManager) dgx.ImageService {
	return &imageService{fm}
}

type imageService struct {
	fileMgr dgx.FileManager
}

func (s *imageService) Upload(ctx context.Context, r io.Reader) (fileID string, err error) {
	if au := dgx.AuthFromContext(ctx); au == nil {
		err = dgx.AuthErrNoAccess
		return
	}

	fileID, err = s.fileMgr.Save(r)
	if err != nil {
		err = errors.New(dgx.ImageErrUpload, err)
		return
	}

	return fileID, nil
}

func (s *imageService) Thumbnail(fileID string, width, height uint) (path string, err error) {
	f, err := s.Image(fileID)
	if err != nil {
		return
	}

	t, err := image.Thumbnail(f, width, height)
	if err != nil {
		err = errors.New(dgx.ImageErrThumbnail, err)
		return
	}

	return t, nil
}

func (s *imageService) Image(fileID string) (path string, err error) {
	path, err = s.fileMgr.Get(fileID)
	if err != nil {
		err = errors.New(dgx.ImageErrNotFound, err)
		return
	}

	return path, nil
}

func (s *imageService) Delete(ctx context.Context, fileID string) error {
	if au := dgx.AuthFromContext(ctx); au == nil {
		return dgx.AuthErrNoAccess
	}

	return s.fileMgr.Delete(fileID)
}
