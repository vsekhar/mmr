package mmr

func intSliceEqual(a, b []int) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil && len(b) == 0:
		return true
	case len(a) == 0 && b == nil:
		return true
	case a == nil || b == nil:
		return false
	case len(a) != len(b):
		return false
	default:
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
}
