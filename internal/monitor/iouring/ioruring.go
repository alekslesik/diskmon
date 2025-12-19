package iouring

// Использование пакета github.com/shirou/gopsutil для мониторинга процессов
import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/process"
)

func MonitorIOURingSyscalls(pid int32) error {
	p, _ := process.NewProcess(pid)
	// Получаем количество системных вызовов (но без детализации по типам)
	// Это косвенный признак активности io_uring

	ioStat, err := p.NetIOCounters(true)
	if err != nil {
		return fmt.Errorf("fail to get net counters of process %w", err)
	}
    
    ticker := time.NewTicker(time.Second * 5)
    
    for {
        for _, s := range ioStat {
            log.Printf("Name %s\n\tBytesSent %d\n\tBytesRecv %d", s.Name, s.BytesSent, s.BytesRecv)
        }
        
        <-ticker.C
    }
}
