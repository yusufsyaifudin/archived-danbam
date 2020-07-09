package app

import (
	"github.com/yusufsyaifudin/danbam/repo/appRepo"
	"github.com/yusufsyaifudin/danbam/server"
)

type handler struct {
	appRepo appRepo.Service
}

func Routes(appRepo appRepo.Service) []*server.Route {
	h := &handler{
		appRepo: appRepo,
	}

	return []*server.Route{
		{
			Path:       "/v1/app",
			Method:     "POST",
			Handler:    h.createApp,
			Middleware: nil,
		},
	}
}
