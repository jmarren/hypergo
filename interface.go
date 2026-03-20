package hypergo

import (
	"context"
	"io"
)

type Component interface {
	// Render the template.
	Render(ctx context.Context, w io.Writer) error
}
