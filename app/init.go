package main

import (
	"os"

	"github.com/mhausenblas/reshifter/pkg/types"
)

func init() {
	// make sure that our working directory exists:
	if _, err := os.Stat(types.DefaultWorkDir); os.IsNotExist(err) {
		_ = os.MkdirAll(types.DefaultWorkDir, 0777)
	}
}
