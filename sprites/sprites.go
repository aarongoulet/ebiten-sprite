package sprites

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"time"
)

type (
	HorizonalAlignment func(x float64, width float64) float64
	VerticalAlignment  func(y float64, height float64) float64
	Alignment          func(x float64, y float64, width float64, height float64) (float64, float64)
)

var (
	Left HorizonalAlignment = func(x float64, width float64) float64 {
		return x
	}

	Center HorizonalAlignment = func(x float64, width float64) float64 {
		return x - (width / 2)
	}

	Right HorizonalAlignment = func(x float64, width float64) float64 {
		return x - width
	}

	Top VerticalAlignment = func(y float64, height float64) float64 {
		return y
	}

	Middle VerticalAlignment = func(y float64, height float64) float64 {
		return y - (height / 2)
	}

	Bottom VerticalAlignment = func(y float64, height float64) float64 {
		return y - height
	}

	TopLeft Alignment = func(x float64, y float64, width float64, height float64) (float64, float64) {
		return Top(y, height), Left(x, width)
	}

	TopCenter Alignment = func(x float64, y float64, width float64, height float64) (float64, float64) {
		return Top(y, height), Center(x, width)
	}

	TopRight Alignment = func(x float64, y float64, width float64, height float64) (float64, float64) {
		return Top(y, height), Right(x, width)
	}

	MiddleLeft Alignment = func(x float64, y float64, width float64, height float64) (float64, float64) {
		return Middle(y, height), Left(x, width)
	}

	MiddleCenter Alignment = func(x float64, y float64, width float64, height float64) (float64, float64) {
		return Middle(y, height), Center(x, width)
	}

	MiddleRight Alignment = func(x float64, y float64, width float64, height float64) (float64, float64) {
		return Middle(y, height), Right(x, width)
	}

	BottomLeft Alignment = func(x float64, y float64, width float64, height float64) (float64, float64) {
		return Bottom(y, height), Left(x, width)
	}

	BottomMiddle Alignment = func(x float64, y float64, width float64, height float64) (float64, float64) {
		return Bottom(y, height), Middle(x, width)
	}

	BottomRight Alignment = func(x float64, y float64, width float64, height float64) (float64, float64) {
		return Bottom(y, height), Right(x, width)
	}
)

type Sprite struct {
	X         float64
	Y         float64
	Animation *Animation
	Scale     float64
	Speed     float64
	Angle     int
	Origin    Alignment
	Frame     int
	Paused    bool
	Visible   bool
	Repeat    bool
	last      time.Time
	options   *ebiten.DrawImageOptions
}

func (s *Sprite) Update() {
	// Don't update paused or non-visible objects.
	if !s.Visible || !s.Paused {
		return
	}

	now := time.Now()
	elapsed := now.Sub(s.last)
	duration := s.Animation.Frames[s.Frame].Duration * time.Duration(s.Speed)

	// Frame change threshold not reached.
	if elapsed < duration {
		return
	}

	// Is this the final frame?
	if s.Frame >= len(s.Animation.Frames)-1 {
		if s.Repeat {
			// Restart the loop.
			s.Frame = 0
		} else {
			// If repeat is disabled, pause the animation after the final frame.
			s.Paused = true
		}
	} else {
		// Not the last frame.  Incrememnt current frame by one.
		s.Frame++
	}

	// Reset the update timer.
	s.last = now
}

func (s *Sprite) Draw(target *ebiten.Image) {
	// Don't draw non-visible objects.
	if !s.Visible {
		return
	}

	// Calculate values.
	frame := s.Animation.Frames[s.Frame]
	w, h := frame.Image.Size()
	drawWidth := float64(w) * s.Scale
	drawHeight := float64(h) * s.Scale
	drawX, drawY := s.Origin(s.X, s.Y, drawWidth, drawHeight)

	// Apply transformations.
	s.options.GeoM.Reset()
	s.options.GeoM.Scale(s.Scale, s.Scale)
	s.options.GeoM.Translate(drawX, drawY)
	s.options.GeoM.Rotate((float64(s.Angle) * math.Pi) / 180)

	target.DrawImage(frame.Image, s.options)
}

func NewSprite(animation *Animation) *Sprite {
	return &Sprite{
		X:         0,
		Y:         0,
		Animation: animation,
		Scale:     1,
		Speed:     1,
		Angle:     0,
		Origin:    TopLeft,
		Frame:     0,
		Paused:    false,
		Visible:   true,
		Repeat:    true,
		last:      time.Now(),
		options:   &ebiten.DrawImageOptions{},
	}
}
