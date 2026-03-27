package transport

import (
	"context"
	"fmt"
	"net/rpc"
	"time"
)

func callWithTimeout(client *rpc.Client, method string, args, reply interface{}, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- client.Call(method, args, reply) }()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("RPC timeout: %s", method)
	}
}
