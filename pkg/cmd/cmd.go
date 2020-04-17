package cmd

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/lunixbochs/vtclean"
)

// Payload is output for a command
type Payload struct {
	Pipe string `json:"pipe" binding:"required"`
	Text string `json:"text" binding:"required"`
}

// Run executes a input command in shell
func Run(command string, sink chan Payload) {
	// Get shell name based on OS
	shell, flag := func() (string, string) {
		if runtime.GOOS == "windows" {
			return "cmd", "/C"
		}
		return "bash", "-c"
	}()
	cmd := exec.Command(shell, flag, command)

	if sink == nil {
		// Run in native mode if no sink
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

	} else {
		// Start command
		stderr, _ := cmd.StderrPipe()
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()
		sink <- Payload{
			Pipe: "echo",
			Text: ">> " + command,
		}

		// Create scan that breaks and emits payload
		scan := func(reader io.Reader, pipe string) {
			scanner := bufio.NewScanner(reader)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				sink <- Payload{
					Pipe: pipe,
					Text: vtclean.Clean(scanner.Text(), false),
				}
			}
		}

		// Scan both stdout and stderr till termination
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			scan(stderr, "stderr")
			wg.Done()
		}()
		scan(stdout, "stdout")
		wg.Wait()

		// Close channel
		close(sink)
	}
}
