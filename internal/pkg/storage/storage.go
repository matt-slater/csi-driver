// Package storage contains interfaces and backends for storing volumes.
package storage

type Storage interface {
	WriteVolume(id, targetPath string, vCtx map[string]string, data []byte) error
	PathForVolume(id string) string
	ListVolumes() []string
	RemoveVolume(id string) error
}

type Volume struct {
	Metadata VolumeMetadata
	Data     []byte
}

type VolumeMetadata struct {
	ID            string
	TargetPath    string
	VolumeContext map[string]string
}
