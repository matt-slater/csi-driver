// Package storage contains interfaces and backends for storing volumes.
package storage

type Storage interface {
	WriteVolume(id string, vCtx map[string]string) (bool, error)
	PathForVolume(id string) string
	ListVolumes() ([]string, error)
	RemoveVolume(id string) error
}
