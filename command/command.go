package command

import (
	"context"
	"os/exec"
	"syscall"
	"time"
)

func Run(ctx context.Context, cmd *exec.Cmd) (<-chan error, error) {
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	wait1 := make(chan error)
	go func() {
		err := cmd.Wait()
		wait1 <- err
		close(wait1)
	}()

	wait2 := make(chan error)
	go func() {
		var wait1Err error
		defer func() {
			wait2 <- wait1Err
			close(wait2)
		}()
		select {
		case wait1Err = <-wait1:
			return
		case <-ctx.Done():
			_ = cmd.Process.Signal(syscall.SIGTERM)
			timer := time.NewTimer(2 * time.Second)
			defer timer.Stop()
			select {
			case wait1Err = <-wait1:
				return
			case <-timer.C:
				_ = cmd.Process.Kill()
				wait1Err = <-wait1
				return
			}
		}
	}()
	return wait2, nil
}
