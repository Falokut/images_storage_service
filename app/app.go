// nolint:containedctx
package app

import (
	"context"

	"github.com/Falokut/go-kit/log"
	"github.com/Falokut/images_storage_service/conf"
	"github.com/pkg/errors"
)

type RunnerFunc func(context.Context) error
type CloserFunc func(context.Context) error

type Application struct {
	logger  log.Logger
	context context.Context
	cfg     config
	runners []RunnerFunc
	closers []CloserFunc
}

type config struct {
	localCfg *conf.LocalConfig
}

type appCtxKey struct{}

func New() *Application {
	localCfg := conf.GetLocalConfig()
	ctx := context.WithValue(context.Background(), appCtxKey{}, localCfg.App)

	level, err := log.ParseLogLevel(localCfg.Log.LogLevel)
	if err != nil {
		panic(errors.WithMessage(err, "parse log level"))
	}
	logCfg := log.Config{
		Loglevel: level,
		Output: log.OutputConfig{
			Console:  localCfg.Log.ConsoleOutput,
			Filepath: localCfg.Log.Filepath,
		},
	}
	logger, err := log.NewFromConfig(logCfg)
	if err != nil {
		panic(errors.WithMessage(err, "logger from config"))
	}

	cfg := config{
		localCfg: localCfg,
	}
	return &Application{
		cfg:     cfg,
		logger:  logger,
		context: ctx,
	}
}

//nolint:ireturn
func (a *Application) GetLogger() log.Logger {
	return a.logger
}

func (a *Application) Context() context.Context {
	return a.context
}

func (a *Application) Config() config {
	return a.cfg
}

func (c config) Local() *conf.LocalConfig {
	return c.localCfg
}

func (a *Application) AddRunners(runners ...RunnerFunc) {
	a.runners = append(a.runners, runners...)
}

func (a *Application) AddClosers(closers ...CloserFunc) {
	a.closers = append(a.closers, closers...)
}

func (a *Application) Run() error {
	for _, runner := range a.runners {
		go func() {
			err := runner(a.context)
			if err != nil {
				a.logger.Fatal(a.context, errors.WithMessage(err, "run runner"))
			}
		}()
	}
	return nil
}

func (a *Application) Shutdown() {
	for _, closer := range a.closers {
		err := closer(a.context)
		if err != nil {
			a.logger.Error(a.context, errors.WithMessage(err, "run closer"))
		}
	}
}
