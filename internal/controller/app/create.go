package app

import (
	"context"

	"github.com/yusufsyaifudin/danbam/server"
)

type requestCreateApp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h handler) createApp(ctx context.Context, req server.Request) server.Response {
	form := &requestCreateApp{}
	err := req.Bind(form)
	if err != nil {
		return nil
	}

	return nil
}
