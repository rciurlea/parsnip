package jobs

import (
	"context"
	"fmt"
)

type A struct {
	Something string
}

func (a A) Run(ctx context.Context) (bool, error) {
	fmt.Println(a)
	return false, nil
}
