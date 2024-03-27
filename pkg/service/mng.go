package service

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
)

func NewManager() (*Manager, error) {
	m := Manager{
		services: make([]Starter, 0),
	}
	return &m, nil
}

type Starter interface {
	// Run controller and block until finish
	Run(ctx context.Context) error
}

type Manager struct {
	services []Starter
}

func (m *Manager) Add(s Starter) {
	m.services = append(m.services, s)
}

func (m *Manager) StartAndWait(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)
	for _, s := range m.services {
		service := s
		group.Go(func() error {
			err := service.Run(ctx)
			if err != nil && !errors.Is(err, context.Canceled) {
				return err
			}
			<-ctx.Done()
			return err
		})
	}
	return group.Wait()
}
