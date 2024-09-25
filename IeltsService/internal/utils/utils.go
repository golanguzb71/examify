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
