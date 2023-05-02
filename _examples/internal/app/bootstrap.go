package app

import (
	"context"
	"github.com/gogoclouds/gogo/_examples/internal/app/admin/model"
	"github.com/gogoclouds/gogo/app"
)

func Run(config string) {
	app.
		New(context.Background(), config).
		Database().
		AutoMigrate(model.Tables).
		HTTP(loadRouter).
		Run()
}