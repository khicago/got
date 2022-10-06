package procast


import (
	"sync"
)

// MultiGo
// - fast return when any error occurred
func MultiGo(handlers ...func() error) error {
	var wg sync.WaitGroup

	wg.Add(len(handlers))
	errCh := make(chan error, len(handlers))

	for i := range handlers {
		t := handlers[i]
		SafeGo(func() {
			defer wg.Done()
			if err := t(); err != nil {
				errCh <- err
			}
		}, func(err error) {
			errCh <- err
		})
	}

	SafeGo(func() {
		wg.Wait()
		close(errCh)
	}, func(err error) {})

	for err := range errCh {
		return err
	}

	return nil
}

