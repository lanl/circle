// This program forms the core of a parallel version of the POSIX
// xargs command.  A command template is provided on the command line,
// and this template is instantiated for each line of input read from
// the standard input device.
package main

import (
	"bufio"
	"github.com/losalamos/circle"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Parse the command line into a template.
	notify := log.New(os.Stderr, "mpixargs: ", 0)
	template := "echo -n '{} '" // Similar to xargs's default
	if len(os.Args) > 1 {
		template = strings.Join(os.Args[1:], " ")
	}
	if !strings.Contains(template, "{}") {
		template = template + " {}"
	}

	// Initialize the circle library.
	rank := circle.Initialize()
	defer circle.Finalize()

	// Create a pair of channels for writing work into the queue
	// and reading work from the queue.
	toQ, fromQ := circle.ChannelBegin()

	// Process 0 enqueues a command per line of standard input.
	if rank == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			toQ <- strings.Replace(template, "{}", scanner.Text(), -1)
		}
		if err := scanner.Err(); err != nil {
			notify.Fatalln(err)
		}
		close(toQ)
	}

	// All processes collaboratively execute commands read from the queue.
	for command := range fromQ {
		cmdState := exec.Command("bash", "-c", command)
		cmdState.Stdout = os.Stdout
		cmdState.Stderr = os.Stderr
		if err := cmdState.Run(); err != nil {
			notify.Printf("WARNING: %s", err)
		}
	}
}
