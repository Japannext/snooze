package patlite

import (
	"fmt"
)

var (
	COLORS     = [5]string{"red", "yellow", "green", "blue", "white"}
	COLOR_HASH = map[string]bool{"red": true, "yellow": true, "green": true, "blue": true, "white": true}
)

var (
	lightStringToByte = map[string]byte{
		"off":    OFF,
		"on":     ON,
		"blink1": BLINK1,
		"blink2": BLINK2,
	}
	lightByteToString = map[byte]string{
		OFF:    "off",
		ON:     "on",
		BLINK1: "blink1",
		BLINK2: "blink2",
	}
	soundStringToByte = map[string]byte{
		"off":   OFF,
		"short": SHORT,
		"long":  LONG,
		"tiny":  TINY,
		"beep":  BEEP,
	}
	soundByteToString = map[byte]string{
		OFF:   "off",
		SHORT: "short",
		LONG:  "long",
		TINY:  "tiny",
		BEEP:  "beep",
	}
)

type State struct {
	Lights map[string]string `json:"lights" yaml:"lights"`
	Sound  string            `json:"sound"  yaml:"sound"`
}

func (s *State) Validate() error {
	for color, light := range s.Lights {
		if ok := COLOR_HASH[color]; !ok {
			return fmt.Errorf("invalid color '%s'", color)
		}
		if _, ok := lightStringToByte[s.Lights[color]]; !ok {
			return fmt.Errorf("invalid light '%s'", light)
		}
	}
	if _, ok := soundStringToByte[s.Sound]; !ok {
		return fmt.Errorf("invalid sound '%s'", s.Sound)
	}
	return nil
}

func newState(b []byte) (*State, error) {
	if len(b) < 6 {
		return nil, fmt.Errorf("unknown response: expected 6 bytes, got %d", len(b))
	}
	s := &State{
		Lights: make(map[string]string),
	}
	for i, color := range COLORS {
		light, ok := lightByteToString[b[i]]
		if !ok {
			return s, fmt.Errorf("unknown byte `%x` for light", b[i])
		}
		s.Lights[color] = light
	}
	var ok bool
	s.Sound, ok = soundByteToString[b[5]]
	if !ok {
		return s, fmt.Errorf("unknown byte '%x' for sound", b[5])
	}

	return s, nil
}

func (s *State) bytes() []byte {
	res := make([]byte, 6)
	for i, color := range COLORS {
		light, ok := s.Lights[color]
		if !ok {
			res[i] = OFF
			continue
		}
		res[i] = lightStringToByte[light]
	}
	res[5] = soundStringToByte[s.Sound]
	return res
}
