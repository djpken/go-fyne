package mobile

import "djpken/go-fyne"

func (d *driver) StartAnimation(a *fyne.Animation) {
	d.animation.Start(a)
}

func (d *driver) StopAnimation(a *fyne.Animation) {
	d.animation.Stop(a)
}
