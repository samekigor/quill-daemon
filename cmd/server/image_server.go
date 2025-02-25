package server

import (
	"context"
	"fmt"

	"github.com/samekigor/quill-daemon/cmd/registry"
	img "github.com/samekigor/quill-daemon/proto/images"
)

type ImgServer struct {
	img.UnimplementedImagesServer
}

func (s *ImgServer) PullImage(ctx context.Context, req *img.PullRequest) (*img.PullStatus, error) {

	re, err := registry.GetRegistryEntryByName(req.Registry)
	if err != nil {
		return &img.PullStatus{
			IsSuccess: false,
			Message:   "Failure with retreiving credentials from file",
		}, err
	}

	if re == nil {
		return &img.PullStatus{
			IsSuccess: false,
			Message:   fmt.Sprintf("Credentials not found. Login to registry %v", req.Registry),
		}, nil
	}

	// re.

	return &img.PullStatus{
		IsSuccess: true,
		Message:   " ",
	}, nil
}
