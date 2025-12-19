package iouring

// Использование пакета github.com/shirou/gopsutil для мониторинга процессов
import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/process"
)

func MonitorIOURingSyscalls() error {
	pids, err := process.Pids()
	if err != nil {
		return fmt.Errorf("fail to get pids %w", err)
	}

	pid := pids[len(pids)/2]

	p, _ := process.NewProcess(pid)
	// Получаем количество системных вызовов (но без детализации по типам)
	// Это косвенный признак активности io_uring

	ioStat, err := p.NetIOCounters(true)
	if err != nil {
		return fmt.Errorf("fail to get net counters of process %w", err)
	}

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for range ticker.C {
		for _, s := range ioStat {
			log.Printf("Name %s\n\tBytesSent %d\n\tBytesRecv %d", s.Name, s.BytesSent, s.BytesRecv)
		}
	}

	// for {
	//     for _, s := range ioStat {
	//         log.Printf("Name %s\n\tBytesSent %d\n\tBytesRecv %d", s.Name, s.BytesSent, s.BytesRecv)
	//     }

	//     <-ticker.C
	// }
	
	return nil
	
}
