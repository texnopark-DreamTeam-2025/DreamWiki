package component

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

type Component interface {
	Name() string

	// Run should run the component and block until the context is cancelled.
	// If component fails, it should return error.
	Run(ctx context.Context) error
}

// tryWriteErrorInChannel writes error only if it would be non-blocking.
func tryWriteErrorInChannel(errCh chan error, err error) {
	select {
	case errCh <- err:
	default:
	}
}

func RunComponents(components ...Component) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, len(components))
	wg := sync.WaitGroup{}
	componentsRunning := int64(0)

	for _, component := range components {
		wg.Add(1)
		go func() {
			defer wg.Done()

			startSequenceNumber := atomic.AddInt64(&componentsRunning, 1)
			log.Printf("Starting component %s (%d/%d)", component.Name(), startSequenceNumber, len(components))

			err := component.Run(ctx)

			componentsStillRunning := atomic.AddInt64(&componentsRunning, -1)
			componentsStopped := int64(len(components)) - componentsStillRunning

			if err != nil && !errors.Is(err, context.Canceled) {
				log.Printf("Component %s stopped (%d/%d), with error: %s", component.Name(), int(componentsStopped), len(components), err.Error())
				tryWriteErrorInChannel(errCh, err)
			} else {
				log.Printf("Component %s stopped (%d/%d)", component.Name(), int(componentsStopped), len(components))
			}
		}()
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case err := <-errCh:
		log.Printf("Shutting down by component failure...")
		cancel()
		wg.Wait()
		return err
	case sig := <-sigCh:
		log.Printf("Received signal %v, shutting down gracefully...", sig)
		cancel()
		wg.Wait()
		return nil
	}
}
