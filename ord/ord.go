package ord

// Prefix is the OP_RETURN prefix for the 1Sat Ordinals inscription protocol
const Prefix string = "ord"

// Ordinal tells wether an inscription is found
type Ordinal struct {
	Inscription bool
}
