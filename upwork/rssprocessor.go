package upwork

import (
	"time"
)

type RSSProcessor struct {
	cache map[string]Entry
}

func NewRSSProcessor() *RSSProcessor {
	return &RSSProcessor{
		cache: make(map[string]Entry),
	}
}

func (rp *RSSProcessor) Check(entries []Entry, skipDuration int) []Entry {
	res := make([]Entry, 0)
	now := time.Now()
	sd := time.Duration(skipDuration) * time.Minute
	for _,v := range entries {
		dif := now.Sub(v.Updated)
		if dif > sd {
			continue
		}
		if _,exist := rp.cache[v.Id]; !exist {
			res = append(res, v)
			rp.cache[v.Id] = v
		}
	}

	for k,v := range rp.cache {
		dif := now.Sub(v.Updated)
		if dif > sd {
			delete(rp.cache,k)
		}
	}

	return res
}