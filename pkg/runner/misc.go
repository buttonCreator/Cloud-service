package runner

import "context"

var (
	_ StartStopper = (*callbackNoopStarer)(nil)
	_ Shutdownable = (*callbackNoopStarer)(nil)
	_ MainService  = (*mainService)(nil)
)

type callbackNoopStarer struct {
	stopCallback func() error
}

type mainService struct {
	startCallback func(ctx context.Context) error
}

func (*callbackNoopStarer) Start(_ context.Context) error {
	return nil
}

func (c *callbackNoopStarer) Shutdown() error {
	if c.stopCallback == nil {
		return nil
	}

	return c.stopCallback()
}

func NewFunctionAsMain(cbStart func(ctx context.Context) error) *mainService {
	return &mainService{startCallback: cbStart}
}

func (m *mainService) Start(ctx context.Context) error {
	if m.startCallback == nil {
		return nil
	}

	return m.startCallback(ctx)
}

func (*mainService) Shutdown() error {
	return nil
}

func (*mainService) SetReady(_ bool) {}
