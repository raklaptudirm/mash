// mash
// https://github.com/raklaptudirm/mash
// Copyright (c) 2021 Rak Laptudirm.
// Licensed under the MIT license.

package builtin

import (
	"fmt"
	"os"

	"github.com/raklaptudirm/mash/commands"
)

// Clear command is used to clear the terminal,
// including scroll-back (for now).
func clear(args []string) error {
	if len(args) > 0 {
		fmt.Fprintln(os.Stderr, "clear: too many arguments")
		return &commands.ExitError{Code: 1}
	}

	// Escape sequence to preserve scroll-back:
	// fmt.Print("\u001b[2J")
	fmt.Print("\u001bc")
	return nil
}
