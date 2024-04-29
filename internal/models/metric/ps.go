package metric

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type PSMetrics struct {
	// метрики пакета gopsutil.
	TotalMemory     Gauge
	FreeMemory      Gauge
	CPUutilization1 Gauge
}

func (ps *PSMetrics) Poll() error {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	ps.TotalMemory = Gauge(vm.Total)
	ps.FreeMemory = Gauge(vm.Free)

	CPUCount, err := cpu.Counts(true)
	if err != nil {
		return err
	}

	ps.CPUutilization1 = Gauge(CPUCount)

	return nil
}
