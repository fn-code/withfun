package service

import (
	"fmt"

	"github.com/fn-code/withfun/testing/internal/storage"
)

type Impl struct {
	Storage storage.Storage
}

func (i *Impl) Sum(a, b int) (int, error) {
	return i.Storage.Add(a, b)
}

func (i *Impl) SumCustom(a, b int, fn func(job string) error) (int, error) {

	return i.Storage.AddCustom(&storage.Custom{A: a, B: b, OnProcess: fn})
}

func reproduceFunc() func(job string) error {
	return func(job string) error {
		fmt.Println("got job : ", job)
		return nil
	}
}
