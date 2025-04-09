package controller

//import "testing"

// Don't do anything, if there is no vehicle in range.
// Don't do anything, if there is no vehicle in range that's charging.
// Return error, if you can't check if vehicle is charging.
// Check that the actions are performed on the vehicle that's charging (when two are in range, but only one is charging).

// Tests for expected amps set:
// surplus    current amps    expected set amps
//  12        0               16
//  0         8               8
// -1         16              5
