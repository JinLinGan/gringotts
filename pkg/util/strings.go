package util

func StringSliceEqual(a, b *[]string) bool {
	if len(*a) != len(*b) {
		return false
	}

	for _, ia := range *a {
		found := false
		for _, ib := range *b {
			if ia == ib {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
