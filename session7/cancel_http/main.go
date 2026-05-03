package main

import (
	"context"
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
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	ch := make(chan string)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			resp, err := makeRequest(ctx, "http://httpbin.org/status/200,200,200,500")
			if err != nil {
				fmt.Println("Error from status goroutine", err)
				cancelFunc()
				return
			}
			if resp.StatusCode == http.StatusInternalServerError {
				fmt.Println("Internal server error,exiting")
				cancelFunc()
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
				fmt.Println("Error from delay goroutine", err)
				cancelFunc()
				return
			}
			select {
			case ch <- "success from delay goroutine" + resp.Header.Get("date"):
			case <-ctx.Done():
			}
		}
	}()
loop:
	for {
		select {
		case msg := <-ch:
			fmt.Println(msg)
		case <-ctx.Done():
			break loop
		}
	}
	wg.Wait()
}
