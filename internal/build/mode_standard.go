//go:build !debug && !release

package build

import "github.com/djpken/go-fyne"

// Mode is the application's build mode.
const Mode = fyne.BuildStandard
