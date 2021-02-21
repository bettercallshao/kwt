package cmd

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
	shell, flag, header := func() (string, string, string) {
		if runtime.GOOS == "windows" {
			return "cmd", "/C", "@echo off\n"
		}
		return "bash", "-c", ""
	}()

	// Write command to a file, don't handle err
	// Add .bat since windows cares about extention
	dir, _ := ioutil.TempDir("", "")
	path := filepath.Join(dir, "kwt.bat")
	ioutil.WriteFile(path, []byte(header+command), 0700)
	defer os.RemoveAll(dir)
	cmd := exec.Command(shell, flag, path)

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
