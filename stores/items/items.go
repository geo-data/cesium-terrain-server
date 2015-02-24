package items

type Item struct {
	value []byte
}

// MarshalBinary implements the encoding.MarshalBinary interface.
func (this *Item) MarshalBinary() ([]byte, error) {
	return this.value, nil
}

// UnmarshalBinary implements the encoding.UnmarshalBinary interface.
func (this *Item) UnmarshalBinary(data []byte) error {
	this.value = data
	return nil
}
