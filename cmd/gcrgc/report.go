package main

type report struct {
	total        int
	totalDeleted int
}

func newReport() *report {
	return &report{0, 0}
}

func (r report) Total() int {
	return r.total
}

func (r report) TotalDeleted() int {
	return r.totalDeleted
}

func (r *report) reportTag(tag Tag) {
	r.total++
	if tag.IsRemoved {
		r.totalDeleted++
	}
}
