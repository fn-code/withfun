package storage

import (
	"errors"
)

var (
	ErrZoro = errors.New("error zero")
)

type Custom struct {
	A, B      int
	OnProcess func(jobname string) error
}


//go:generate mockgen -destination=mock/storage.go -package=mock github.com/fn-code/withfun/testing/internal/storage Storage
type Storage interface {
	Add(int, int) (int, error)
	AddCustom(custom *Custom) (int, error)
}

type Calc struct {
}

func NewCalcFactory() Storage {
	return &Calc{}
}

func (c *Calc) Add(i int, i2 int) (int, error) {
	if i == 0 && i2 == 0 {
		return 0, ErrZoro
	}

	return i + i2, nil
}

func (c *Calc) AddCustom(custom *Custom) (int, error) {
	if custom.A == 0 && custom.B == 0 {
		return 0, ErrZoro
	}

	//if custom.OnProcess != nil {
	//
	//	err := custom.OnProcess("storage job")
	//	if err != nil {
	//		log.Println("-------------------------------- got error")
	//		return 0, ErrZoro
	//	}
	//}

	return custom.A + custom.B, nil
}
