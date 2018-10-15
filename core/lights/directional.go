package lights

import (
	"../math"
)

type Directional struct {
	Light
}

func NewDirectional(color math.Color, intensity float32) *Directional {
	light := Directional{
		Light: *NewLight(color, intensity),
	}

	return &light
}
