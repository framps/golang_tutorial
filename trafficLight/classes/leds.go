package classes

var T1LEDs = LEDs{[...]int{11, 12, 13}}
var T2LEDs = LEDs{[...]int{21, 22, 23}}

// LEDs - LED pin numbers for lights of one traffic light
type LEDs struct {
	pin [3]int // red, yellow, green
}
