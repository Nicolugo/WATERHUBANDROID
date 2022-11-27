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
	//