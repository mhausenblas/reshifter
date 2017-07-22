package main

import (
	"os"

	"github.com/mhausenblas/reshifter/app/handler"
	"github.com/mhausenblas/reshifter/pkg/types"
)

var (
	releaseVersion string
)

func init() {
	handler.ReleaseVersion = releaseVersion
	// make sure that our working directory exists:
	if _, err := os.Stat(types.DefaultWorkDir); os.IsNotExist(err) {
		_ = os.MkdirAll(types.DefaultWorkDir, 0777)
	}
}
