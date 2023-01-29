package os

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"os/exec"
	"strings"
	"time"
)

func GetCmdResult(cmdName string, args ...string) ([]byte, error) {
	timeout := 60 * 5
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmdName, args...)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(stdout))
	}
	lastline := strings.Split(string(stdout), "\n")
	return []byte(lastline[len(lastline)-2]), nil
}
