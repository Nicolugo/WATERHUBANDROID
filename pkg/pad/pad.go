package pad

// This is the matrix layout for register $FF00:
//              P14      P15
//               |        |
//  P10-------O-Right----O-A
//               |        |
//  P11-------O-Left-----O-B
//               |        |
//  P12-------O-Up-------O-Select
//              |         |
//  P13-------O-Down-----O-Start
//

type Pad struct {
	// Bit 7 - Not used
	// Bit 6 - Not used
	// Bit 5 - P15 out port
	// Bit 4 - P14 out port
	// Bit 3 - P13 in port
	// Bit 2 - P12 in port
	// Bit 1 - P11 in port
	// Bit 0 - P10 in port
	reg   byte
	state Button
}

type Button byte

const (
	// A is the A button on the GameBoy.
	A Button = 