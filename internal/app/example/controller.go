package example

import (
	"fmt"
)

type Controller struct {
	Service *Service
}

func NewController() *Controller {
	return &Controller{
		Service: NewService(),
	}
}