package worker

import (
	"context"
)

type Worker interface {
	Run(context.Context, int) error
}
