package command

import (
	"context"
	"os/exec"
	"syscall"
	"time"
)

func Run(ctx context.Context, cmd *exec.Cmd) (<-chan struct{}, error) {
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	wait := make(chan struct{})
	go func() {
		defer close(wait)
		cmd.Wait()
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		select {
		case <-wait:
			return
		case <-ctx.Done():
			_ = cmd.Process.Signal(syscall.SIGTERM)
			timer := time.NewTimer(2 * time.Second)
			defer timer.Stop()
			select {
			case <-wait:
				return
			case <-timer.C:
				_ = cmd.Process.Kill()
				<-wait
				return
			}
		}
	}()
	return done, nil
}
