package uid

import (
	"context"
	"errors"
	"github.com/sony/sonyflake"
	"time"
)

type GeneratorConfig struct {
	StartTime      time.Time
	MachineID      func() (uint16, error)
	CheckMachineID func(uint16) bool
}

type Generator interface {
	NextID() (uint64, error)
}

type generator struct {
	sf *sonyflake.Sonyflake
}

func NewGenerator(ctx context.Context) (Generator, error) {
	var st sonyflake.Settings // TODO : settings ref: https://github.com/sony/sonyflake
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		return nil, errors.New("cannot create Sonyflake generator")
	}
	return &generator{sf: sf}, nil
}

func (g *generator) NextID() (uint64, error) {
	// TODO: twist bits?
	return g.sf.NextID()
}
