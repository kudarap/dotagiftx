package service

import (
	"context"
	"io"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
	"github.com/kudarap/dotagiftx/gokit/file/image"
)

// NewAuth returns a new Image service.
func NewImage(fm core.FileManager) core.ImageService {
	return &imageService{fm}
}

type imageService struct {
	fileMgr core.FileManager
}

func (s *imageService) Upload(ctx context.Context, r io.Reader) (fileID string, err error) {
	if au := core.AuthFromContext(ctx); au == nil {
		err = core.AuthErrNoAccess
		return
	}

	fileID, err = s.fileMgr.Save(r)
	if err != nil {
		err = errors.New(core.ImageErrUpload, err)
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
		err = errors.New(core.ImageErrThumbnail, err)
		return
	}

	return t, nil
}

func (s *imageService) Image(fileID string) (path string, err error) {
	path, err = s.fileMgr.Get(fileID)
	if err != nil {
		err = errors.New(core.ImageErrNotFound, err)
		return
	}

	return path, nil
}

func (s *imageService) Delete(ctx context.Context, fileID string) error {
	if au := core.AuthFromContext(ctx); au == nil {
		return core.AuthErrNoAccess
	}

	return s.fileMgr.Delete(fileID)
}
