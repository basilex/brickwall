// BrickBite platform
// CopyRight (C) 2021 by BrickBite GmbH. All Rights Reserved.
package main

import (
	"context"
	"log"

	"brickwall/cmd"
	"brickwall/internal/common"
)

var (
	Version, Staging, Githash, Gobuild, Compile string
)

func main() {
	ctx := context.Background()

	ctx = context.WithValue(ctx,
		common.KeyMetadata, &common.Metadata{
			Version: Version, Staging: Staging, Githash: Githash, Gobuild: Gobuild, Compile: Compile,
		},
	)
	if err := cmd.Bootstrap(ctx); err != nil {
		log.Fatalf("error: %s", err)
	}
}
