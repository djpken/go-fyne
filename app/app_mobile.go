//go:build !ci && (android || ios || mobile)

package app

import (
	"djpken/go-fyne"
	internalapp "djpken/go-fyne/internal/app"
	"djpken/go-fyne/internal/driver/mobile"
)

// NewWithID returns a new app instance using the appropriate runtime driver.
// The ID string should be globally unique to this app.
func NewWithID(id string) fyne.App {
	d := mobile.NewGoMobileDriver()
	a := newAppWithDriver(d, id)
	d.(mobile.ConfiguredDriver).SetOnConfigurationChanged(func(c *mobile.Configuration) {
		internalapp.SystemTheme = c.SystemTheme

		a.Settings().(*settings).setupTheme()
	})
	return a
}
