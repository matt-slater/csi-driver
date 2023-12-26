package storage

import (
	"fmt"
	"sync"
)

type InMemory struct {
	files map[string]Volume
	lock  sync.Mutex
}

func NewInMemory() *InMemory {
	return &InMemory{
		files: make(map[string]Volume),
	}
}

func (m *InMemory) WriteVolume(id, targetPath string, vCtx map[string]string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.files[id] = Volume{
		Metadata: VolumeMetadata{
			ID:            id,
			TargetPath:    targetPath,
			VolumeContext: vCtx,
		},
	}
	return nil
}

func (m *InMemory) PathForVolume(id string) string {
	m.lock.Lock()
	defer m.lock.Unlock()
	return id
}

func (m *InMemory) ReadFile(id string) ([]byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	volume, ok := m.files[id]
	if !ok {
		return nil, fmt.Errorf("failed to find volume id %s in memory-backed store", id)
	}

	cpy := []byte{}

	copy(cpy, volume.Data)

	return cpy, nil
}

func (m *InMemory) ListVolumes() ([]string, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	vols := make([]string, len(m.files))

	i := 0
	for k := range m.files {
		vols[i] = k
		i++
	}

	return vols, nil
}

func (m *InMemory) RemoveVolume(id string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.files, id)
	return nil
}
