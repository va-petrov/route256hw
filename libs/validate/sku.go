package validate

func SKU(sku uint32) error {
	if sku == 0 {
		return ErrEmptySKU
	}
	return nil
}
