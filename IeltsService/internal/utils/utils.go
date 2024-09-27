package utils

func OffSetGenerator(page, size *int32) int {
	if page == nil || *page < 1 {
		p := int32(1)
		page = &p
	}
	if size == nil || *size < 1 {
		s := int32(10)
		size = &s
	}

	return int(*size * (*page - 1))
}

func CalculateBandScore(correctCount int) float64 {
	switch {
	case correctCount >= 39:
		return 9.0
	case correctCount >= 37:
		return 8.5
	case correctCount >= 35:
		return 8.0
	case correctCount >= 32:
		return 7.5
	case correctCount >= 30:
		return 7.0
	case correctCount >= 26:
		return 6.5
	case correctCount >= 23:
		return 6.0
	case correctCount >= 18:
		return 5.5
	case correctCount >= 16:
		return 5.0
	case correctCount >= 13:
		return 4.5
	case correctCount >= 10:
		return 4.0
	case correctCount >= 6:
		return 3.5
	case correctCount >= 4:
		return 3.0
	case correctCount >= 1:
		return 2.5
	default:
		return 2.0
	}
}
