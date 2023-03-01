package validate

func OrderId(orderID int64) error {
	if orderID == 0 {
		return ErrEmptyOrderID
	}
	return nil
}
