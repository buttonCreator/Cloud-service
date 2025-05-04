package runner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const startTimeout = 30 * time.Second

type (
	StartStopper interface {
		Shutdownable
		Start(ctx context.Context) error
	}
	Shutdownable interface {
		Shutdown() error
	}
	Logger interface {
		Info(msg string, keyAndValues ...any)
		Debug(msg string, keyAndValues ...any)
		Sync() error
	}
	SetReady interface {
		SetReady(valueBool bool)
	}
)

type MainService interface {
	StartStopper
	SetReady
}

type runner struct {
	main      MainService
	auxiliary []StartStopper
	logger    Logger
}

func New(logger Logger, main MainService, auxiliary ...StartStopper) *runner {
	return &runner{
		main:      main,
		auxiliary: auxiliary,
		logger:    logger,
	}
}

func (r *runner) startAuxiliary(ctx context.Context) error {
	for _, auxiliary := range r.auxiliary {
		r.logger.Info("starting service", "service", fmt.Sprintf("%T", auxiliary))
		if err := auxiliary.Start(ctx); err != nil {
			return fmt.Errorf("startAuxiliary service '%T': %w", auxiliary, err)
		}
	}

	return nil
}

func (r *runner) stopAuxiliary() error {
	for i := len(r.auxiliary) - 1; i >= 0; i-- {
		x := r.auxiliary[i]
		r.logger.Info("stopAuxiliary service...", "service", fmt.Sprintf("%T", x))
		if err := x.Shutdown(); err != nil {
			return fmt.Errorf("stopAuxiliary service '%T': %w", x, err)
		}
	}

	return nil
}

func (r *runner) RunUtilsSignalExit() error {
	startCtx, cancel := context.WithTimeout(context.Background(), startTimeout)
	defer cancel()

	if err := r.startAuxiliary(startCtx); err != nil {
		return fmt.Errorf("startAuxiliary: %w", err)
	}

	r.logger.Info("starting main", "main", fmt.Sprintf("%T", r.main))
	if err := r.main.Start(startCtx); err != nil {
		return fmt.Errorf("startAuxiliary service '%T': %w", r.main, err)
	}
	r.main.SetReady(true)

	r.logger.Debug("listen signal to shutdown...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	r.logger.Info("received signal. Shutdown...", "signal", <-quit)

	r.main.SetReady(false)

	if err := r.stopAuxiliary(); err != nil {
		return fmt.Errorf("stopAuxiliary: %w", err)
	}

	r.logger.Info("shutdown successfully")

	return nil
}
