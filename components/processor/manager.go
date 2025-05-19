package processor

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

// A manager in order to reload the config
type Manager struct {
	mu sync.Mutex

	processor *Processor
}

func NewManager() *Manager {
	manager := &Manager{}

	if err := manager.Reload(); err != nil {
		log.Fatalf("failed to initialize config: %s", err)
	}

	return manager
}

func (manager *Manager) Reload() error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %s", err)
	}

	newProcessor, err := NewProcessor(cfg)
	if err != nil {
		return fmt.Errorf("while initializing main process: %s", err)
	}
	manager.processor = newProcessor

	return nil
}

func (manager *Manager) GetProcessor() *Processor {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	return manager.processor
}
