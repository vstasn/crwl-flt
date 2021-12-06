package crwlrs

func ContainsNum(values []int64, value int64) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}
