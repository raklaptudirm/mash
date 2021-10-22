// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

// mash is a simple shell written in go.
// Features:
// - cd command
// - exit command
// - run executable files

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/raklaptudirm/mash/vm"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Catch ctrl+c and SIGTERM events so as not to
	// interrupt the shell input, unlike normal
	// processes.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	signal.Notify(interrupt, syscall.SIGTERM)

	// Infinite command loop:
	// - Print prompt
	// - Read input
	// - Parse and Run input
	for {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		// The shell prompt, currently fixed as:
		//
		//     <current working directory>
		//     ψ | <- cursor
		fmt.Printf("\u001b[32m%v\u001b[0m\nψ ", cwd)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		vm.Run(input)
	}
}
