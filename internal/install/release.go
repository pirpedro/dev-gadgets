package install

import (
	"context"
	"fmt"

	"github.com/pirpedro/dev-gadgets/internal/catalog"
)

func runRelease(ctx context.Context, it catalog.Item) error {
	// TODO: download it.Strategy.Release["url"], extract, place bin in ~/.local/bin (use map["bin"] for name)
	return fmt.Errorf("release strategy not implemented: %s", it.ID)
}
