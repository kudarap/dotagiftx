package service

import (
	"context"
	"io"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/file"
	"github.com/kudarap/dotagiftx/xerrors"
)

// NewImage returns a new Image service.
func NewImage(fm dotagiftx.FileManager) dotagiftx.ImageService {
	return &imageService{fm}
}

type imageService struct {
	fileMgr dotagiftx.FileManager
}

func (s *imageService) Upload(ctx context.Context, r io.Reader) (fileID string, err error) {
	if au := dotagiftx.AuthFromContext(ctx); au == nil {
		err = dotagiftx.AuthErrNoAccess
		return
	}

	fileID, err = s.fileMgr.Save(r)
	if err != nil {
		err = xerrors.New(dotagiftx.ImageErrUpload, err)
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
		err = xerrors.New(dotagiftx.ImageErrThumbnail, err)
		return
	}

	return t, nil
}

func (s *imageService) Image(fileID string) (path string, err error) {
	path, err = s.fileMgr.Get(fileID)
	if err != nil {
		err = xerrors.New(dotagiftx.ImageErrNotFound, err)
		return
	}

	return path, nil
}

func (s *imageService) Delete(ctx context.Context, fileID string) error {
	if au := dotagiftx.AuthFromContext(ctx); au == nil {
		return dotagiftx.AuthErrNoAccess
	}

	return s.fileMgr.Delete(fileID)
}
