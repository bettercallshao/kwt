package cmd

import (
	"bufio"
	"github.com/lunixbochs/vtclean"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
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

	// Start command
	cmd := exec.Command(shell, flag, command)
	cmd.Stdin = os.Stdin
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
