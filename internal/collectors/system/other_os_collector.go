//go:build !linux

package system

import (
	"errors"
	"github.com/rsdenck/nux/internal/core/ports"
)

type OtherOSSystemMonitor struct{}

func NewSystemMonitor() ports.SystemMonitor {
	return &OtherOSSystemMonitor{}
}

func (m *OtherOSSystemMonitor) GetMetrics() (ports.SystemMetrics, error) {
	return ports.SystemMetrics{}, errors.New("system monitoring not supported on this OS")
}
