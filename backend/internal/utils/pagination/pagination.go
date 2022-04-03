package pagination

func CalculateOffset(offset, page, limit int) int {
	if offset == 0 && page > 0 {
		offset = (page - 1) * limit
	}

	return offset
}

func CalculatePages(offset, limit, total int) (page, pages int) {
	page, pages = 1, 1

	pages = (total % limit)

	if (total%limit) != 0 || pages == 0 {
		pages++
	}

	page = (offset % limit)

	if (offset%limit) != 0 || page == 0 {
		page++
	}

	return page, pages
}
