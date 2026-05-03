package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func makeRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func main() {
	ctx, cancelFunc := context.WithCancelCause(context.Background())
	defer cancelFunc(nil)

	ch := make(chan string)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			resp, err := makeRequest(ctx, "http://httpbin.org/status/200,200,200,500")
			if err != nil {
				cancelFunc(fmt.Errorf("Error from status goroutine: %w", err))
				return
			}
			if resp.StatusCode == http.StatusInternalServerError {
				cancelFunc(errors.New("Internal server error,exiting"))
				return
			}
			select {
			case ch <- "success from status goroutine":
			case <-ctx.Done():
			}
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		defer wg.Done()
		for {
			resp, err := makeRequest(ctx, "http://httpbin.org/delay/1")
			if err != nil {
				fmt.Println("in delay goroutine:", err)
				cancelFunc(fmt.Errorf("Error from delay goroutine: %w", err))
				return
			}
			select {
			case ch <- "success from delay goroutine " + resp.Header.Get("date"):
			case <-ctx.Done():
			}
		}
	}()
loop:
	for {
		select {
		case msg := <-ch:
			fmt.Println("In main loop:", msg)
		case <-ctx.Done():
			fmt.Println("In main loop: cancelled with error:", context.Cause(ctx))
			break loop
		}
	}
	wg.Wait()
	fmt.Println("Context cancellation cause:", context.Cause(ctx))
}
