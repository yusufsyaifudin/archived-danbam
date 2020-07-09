package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/yusufsyaifudin/danbam/proto"
)

type rpcService struct {
	repo Repo
}

func (r rpcService) CreateApp(ctx context.Context, req *proto.CreateAppRequest) (*proto.CreateAppResponse, error) {
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}

	return &proto.CreateAppResponse{
		App: &proto.App{
			Id:   name,
			Name: name,
		},
	}, nil
}

func NewGRpc(repo Repo) (proto.AppServiceServer, error) {
	if repo == nil {
		panic("repo on app.NewGRpc is nil")
	}

	return &rpcService{
		repo: repo,
	}, nil
}
