package provider

import (
	"errors"
	"sync"
)

// Probably I'll change context opeations with all the providers
// to the approach such as below, which is much safer and quicker.
type IProviderContainer interface {
	Register(string, any)
	Get(string) (any, error)
}

type ProviderContainer struct {
	mu        sync.RWMutex
	providers map[string]any
}

func NewProviderContainer() IProviderContainer {
	return &ProviderContainer{
		providers: make(map[string]any),
	}
}

func (pc *ProviderContainer) Register(key string, provider any) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.providers[key] = provider
}

func (pc *ProviderContainer) Get(key string) (any, error) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	provider, ok := pc.providers[key]
	if !ok {
		return nil, errors.New("provider not found")
	}
	return provider, nil
}

func GetAs[T any](pc *ProviderContainer, key string) (T, error) {
	var zero T

	provider, err := pc.Get(key)
	if err != nil {
		return zero, err
	}
	typedProvider, ok := provider.(T)
	if !ok {
		return zero, errors.New("invalid provider type")
	}
	return typedProvider, nil
}
