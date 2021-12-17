package main

import "time"

// Lifx describes a lightbulb
type Lifx struct {
	ID               string    `json:"id"`
	UUID             string    `json:"uuid"`
	Label            string    `json:"label"`
	Power            string    `json:"power"`
	Connected        bool      `json:"connected"`
	Brightness       float64   `json:"brightness"`
	Color            Color     `json:"color"`
	Effect           string    `json:"effect"`
	Group            Group     `json:"group"`
	Location         Group     `json:"location,omitempty"`
	Product          Product   `json:"product,omitempty"`
	LastSeen         time.Time `json:"last_seen"`
	SecondsSinceSeen int       `json:"seconds_since_seen"`
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

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	Name         string       `json:"name,omitempty"`
	Identifier   string       `json:"identifier,omitempty"`
	Company      string       `json:"company,omitempty"`
	Capabilities Capabilities `json:"capabilities,omitempty"`
	ProductID    int          `json:"product_id,omitempty"`
	VendorID     int          `json:"vendor_id,omitempty"`
}

type Capabilities struct {
	HasColor             bool `json:"has_color,omitempty"`
	HasVariableColorTemp bool `json:"has_variable_color_temp,omitempty"`
	HasIR                bool `json:"has_ir,omitempty"`
	HasChain             bool `json:"has_chain,omitempty"`
	HasMultizone         bool `json:"has_multizone,omitempty"`
	MinKelvin            int  `json:"min_kelvin,omitempty"`
	MaxKelvin            int  `json:"max_kelvin,omitempty"`
}
