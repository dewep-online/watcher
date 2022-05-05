package exc

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/deweppro/go-errors"
)

func Call(ctx context.Context, dir, command string) error {
	cmd := exec.CommandContext(ctx, "/bin/sh", "-xec", fmt.Sprintln(command, " <&-"))
	cmd.Dir = dir
	cmd.Env = os.Environ()

	outp, err0 := cmd.StdoutPipe()
	errp, err1 := cmd.StderrPipe()
	if err := errors.Wrap(err0, err1); err != nil {
		return err
	}

	mr := io.MultiReader(outp, errp)

	go func() {
		scanner := bufio.NewScanner(mr)
		for scanner.Scan() {
			fmt.Println("\t" + scanner.Text())
		}
	}()

	return cmd.Run()
}
