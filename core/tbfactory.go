package core

import "log"

type TBFactory struct {
}

type FactoryCreate struct {
	log.Logger
	Controllers []interface{}
	Providers   []interface{}
}

type TBFactoryInterface interface {
	Create() *TBApp
}

// Create It returns a pointer to a new TBApp
func (tbf *TBFactory) Create(fc *FactoryCreate) *TBApp {
	return &TBApp{}
}
