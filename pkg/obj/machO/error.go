package machO

// InvalidMachOMagicError is returned when the read magic value does not match the expected Mach-O value
type InvalidMachOMagicError struct {
	InvalidMagic string
}

func (e *InvalidMachOMagicError) Error() string {
	return "invalid Mach-O magic: " + e.InvalidMagic
}
