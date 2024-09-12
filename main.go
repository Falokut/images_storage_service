//nolint:importShadow
package main

import (
	"fmt"

	"github.com/Falokut/go-kit/app"
	"github.com/Falokut/go-kit/shutdown"
	"github.com/Falokut/images_storage_service/assembly"
)

//	@title			images_storage_service
//	@version		1.0.0
//	@description	Сервис для хранения изображений
//	@BasePath		/api/images-storage-service

//go:generate swag init --parseDependency
//go:generate rm -f docs/swagger.json docs/docs.go
func main() {
	app, err := app.New()
	if err != nil {
		fmt.Println("shutdown: error while creating app ", err)
		return
	}
	logger := app.GetLogger()
	assembly, err := assembly.New(app.Context(), logger)
	if err != nil {
		logger.Fatal(app.Context(), err)
	}
	app.AddRunners(assembly.Runners()...)
	app.AddClosers(assembly.Closers()...)

	err = app.Run()
	if err != nil {
		app.Shutdown()
		logger.Fatal(app.Context(), err)
	}

	shutdown.On(func() {
		logger.Info(app.Context(), "starting shutdown")
		app.Shutdown()
		logger.Info(app.Context(), "shutdown completed")
	})
}
