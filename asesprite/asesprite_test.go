package asesprite

import (
	"github.com/aarongoulet/ebiten-sprite/v2/sprites"
	"image"
	"reflect"
	"testing"
)

var testJSON = []byte(`{
  "frames": [
  {
    "filename": "test_0",
    "frame": { "x": 0, "y": 0, "w": 10, "h": 10 },
    "rotated": false,
    "trimmed": false,
    "spriteSourceSize": { "x": 0, "y": 0,"w": 10, "h": 10 },
    "sourceSize": { "w": 10, "h": 10 },
    "duration": 100
  },
  {
    "filename": "test_1",
    "frame": { "x": 10, "y": 0, "w": 10, "h": 10 },
    "rotated": false,
    "trimmed": false,
    "spriteSourceSize": { "x": 10, "y": 0,"w": 10, "h": 10 },
    "sourceSize": { "w": 10, "h": 10 },
    "duration": 100
  }
  ],
  "meta": {
    "app": "https://www.aseprite.org/",
    "version": "1.2.3",
    "image": "test.png",
    "format": "I8",
    "size": { "w": 20, "h": 10 },
    "scale": "1",
    "frameTags": [
      { "name": "test1", "from": 0, "to": 1, "direction": "forward" },
      { "name": "test2", "from": 1, "to": 1, "direction": "forward" }
    ],
    "layers": [
      { "name": "base", "opacity": 255, "blendMode": "normal" }
    ],
    "slices": [
    ]
  }
}`)

var testImage = image.NewRGBA(image.Rect(0, 0, 20, 10))

func createTestSpriteSheetLoader() (sprites.SpriteSheetLoader, error) {
	return NewSpriteSheetLoader(testJSON, testImage)
}

func TestNewSpriteSheetLoader(t *testing.T) {
	_, err := createTestSpriteSheetLoader()

	// Check for errors.
	if err != nil {
		t.Fatalf("SpriteSheetLoader creation failed: %s", err)
	}
}

func TestAnimation(t *testing.T) {
	loader, _ := createTestSpriteSheetLoader()

	// Test loading an existing animation.
	animation, err := loader.Animation("test1")

	if err != nil {
		t.Fatalf("failed to load test animation: %s", err)
	}

	if len(animation.Frames) != 2 {
		t.Fatalf("unexpected number of frames (expected 2, got %d)", len(animation.Frames))
	}
}

func TestAnimationNotFound(t *testing.T) {
	loader, _ := createTestSpriteSheetLoader()

	// Test loading a missing animation.
	_, err := loader.Animation("missing")
	expected := AnimationNotFoundError("missing")

	if !reflect.DeepEqual(err, expected) {
		t.Fatalf("expected AnimationNotFoundError, got: %s", err)
	}
}

func TestAllAnimations(t *testing.T) {
	loader, _ := createTestSpriteSheetLoader()

	animations, err := loader.AllAnimations()

	if err != nil {
		t.Fatalf("failed to load animations: %s", err)
	}

	if len(animations) != 2 {
		t.Fatalf("unexpected number of animations (expected 2, got %d)", len(animations))
	}
}
