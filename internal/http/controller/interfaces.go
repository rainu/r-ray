package controller

import (
	"github.com/rainu/r-ray/internal/processor"
)

type Processor interface {
	Process(input processor.Input) (processor.Output, error)
}
