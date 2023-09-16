package models

import "sync"

type CancelManager struct {
	CancelCommands map[string]chan struct{}
	CancelMutex sync.Mutex
}

func NewCancelManager() *CancelManager {
	return &CancelManager {
		CancelCommands: make(map[string]chan struct{}),
	}
}

func (cm *CancelManager) AddCommand(key string, cancel chan struct{}) {
	cm.CancelMutex.Lock()
	defer cm.CancelMutex.Unlock()
	cm.CancelCommands[key] = cancel
}

func (cm *CancelManager) RemoveCommand(key string) {
	cm.CancelMutex.Lock()
	defer cm.CancelMutex.Unlock()
	if cancel, found := cm.CancelCommands[key]; found {
		close(cancel)
		delete(cm.CancelCommands, key)
	}
} 