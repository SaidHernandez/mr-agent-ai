package theme

import (
	"context"
	"sync"
	"time"
)

// Widget is implemented by every data source.
type Widget interface {
	Name() string
	Fetch(ctx context.Context) Result
}

// Result carries the output of a widget fetch.
// Err is non-nil on failure; Value is the display string.
type Result struct {
	WidgetName string
	Value      string
	Err        error
}

// RunAll executes all widgets concurrently and returns results in registration
// order. Enforces a hard 1.4s deadline so the total render time stays under 1.5s.
func RunAll(widgets []Widget) []Result {
	ctx, cancel := context.WithTimeout(context.Background(), 1400*time.Millisecond)
	defer cancel()

	results := make([]Result, len(widgets))
	var wg sync.WaitGroup
	wg.Add(len(widgets))

	for i, w := range widgets {
		i, w := i, w
		go func() {
			defer wg.Done()
			r := w.Fetch(ctx)
			r.WidgetName = w.Name()
			results[i] = r
		}()
	}

	wg.Wait()
	return results
}
