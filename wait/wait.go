package wait

import (
	"context"
	"net"
	"time"
)

func TCP(ctx context.Context, addr string) error {
	var d net.Dialer
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			conn, err := d.DialContext(ctx, "tcp", addr)
			if err != nil {
				continue
			}
			conn.Close()
			return nil
		}
	}
}
