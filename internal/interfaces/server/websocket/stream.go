package websocket

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrInvalidStreamInterval = errors.New(
	"websocket stream interval must be positive",
)

type DataProvider func(ctx context.Context) (any, error)

func Stream(ctx context.Context, connection *Connection, interval time.Duration, provider DataProvider) error {
	if connection == nil {
		return ErrNilConnection
	}

	if interval <= 0 {
		return ErrInvalidStreamInterval
	}

	if provider == nil {
		return errors.New("websocket data provider is nil")
	}

	if err := sendProviderResult(
		ctx,
		connection,
		provider,
	); err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			if err := sendProviderResult(
				ctx,
				connection,
				provider,
			); err != nil {
				return err
			}
		}
	}
}

func sendProviderResult(ctx context.Context, connection *Connection, provider DataProvider) error {
	data, err := provider(ctx)
	if err != nil {
		return fmt.Errorf(
			"provide websocket stream data: %w",
			err,
		)
	}

	if err := connection.WriteJSON(ctx, data); err != nil {
		return fmt.Errorf(
			"write websocket stream data: %w",
			err,
		)
	}

	return nil
}
