package main

// Lifx describes a lightbulb
type Lifx struct {
	ID         string  `json:"id"`
	Label      string  `json:"label"`
	Power      string  `json:"power"`
	Brightness float64 `json:"brightness"`
	Color      Color   `json:"color"`
}

type Color struct {
	Hue        float64 `json:"hue"`
	Kelvin     float64 `json:"kelvin"`
	Saturation float64 `json:"saturation"`
}

type State struct {
	Selector   string  `json:"selector"`
	Brightness float64 `json:"brightness,omitempty"`
	Power      string  `json:"power,omitempty"`
	Color      string  `json:"color,omitempty"`
	Duration   float64 `json:"duration,omitempty"`
}

type States struct {
	States []State `json:"states"`
}
