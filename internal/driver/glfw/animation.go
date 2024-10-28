package glfw

import "github.com/djpken/go-fyne"

func (d *gLDriver) StartAnimation(a *fyne.Animation) {
	d.animation.Start(a)
}

func (d *gLDriver) StopAnimation(a *fyne.Animation) {
	d.animation.Stop(a)
}
