package io

import (
	"fmt"

	"github.com/shirou/gopsutil/process"
)

type PIDInfo struct {
	PIDS map[int32]PID
}

func New() *PIDInfo {
	p := &PIDInfo{}

	return p
}

func (p *PIDInfo) ReadAllIOBytes() (uint64, uint64, uint64, error) {
	var r, w, c uint64

	p.update()
	
	for _, p := range p.PIDS {
		r =  p.GetRBytes()
	}

	return r, w, c, nil
}

func (p *PIDInfo) update() error {
	p.PIDS = make(map[int32]PID)

	pids, err := process.Pids()
	if err != nil {
		return fmt.Errorf("fail to get pids %w", err)
	}

	for _, pd := range pids {
		p.PIDS[pd] = PID{ID: pd}
	}

	return nil
}

type PID struct {
	ID       int32
	Path     string
	PidBytes PidBytes
}

type PidBytes struct {
	Rbytes uint64
	Wbytes uint64
	Cbytes uint64
}

func (p *PID) GetRBytes() uint64 {

	return 0
}
