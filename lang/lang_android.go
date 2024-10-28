package lang

import (
	"github.com/djpken/go-fyne/internal/driver/mobile/app"

	"github.com/jeandeaual/go-locale"
)

func init() {
	locale.SetRunOnJVM(app.RunOnJVM)
}
