package app

import (
	"context"
)

type App struct {
	serviceProvider *serviceProvider
}

func NewApp(ctx context.Context) *App {
	a := &App{}
	a.serviceProvider = NewServiceProvider(ctx)
	return a
}

func (a *App) Run() {
	a.serviceProvider.Start()
}

