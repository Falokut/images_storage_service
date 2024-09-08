//nolint:importShadow
package main

import (
	"github.com/Falokut/images_storage_service/app"
	"github.com/Falokut/images_storage_service/assembly"
	"github.com/Falokut/images_storage_service/shutdown"
)

//	@title			images_storage_service
//	@version		1.0.0
//	@description	Сервис для хранения изображений
//	@BasePath		/api/images-storage

//go:generate swag init --parseDependency
//go:generate rm -f docs/swagger.json docs/docs.go
func main() {
	app := app.New()
	logger := app.GetLogger()

	assembly, err := assembly.New(app.Context(), logger, app.Config().Local())
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
