package gcrgc

import "github.com/graillus/gcrgc/pkg/gcloud"

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

func (r *report) reportImage(img gcloud.Image) {
	r.total++
	if img.IsRemoved {
		r.totalDeleted++
	}
}
