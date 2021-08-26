package upwork


type RSSProcessor struct {
	cache map[string]UWEntry
}

func NewRSSProcessor() *RSSProcessor {
	return &RSSProcessor{}
}

func (rp *RSSProcessor) Check(entries []UWEntry) []UWEntry {
	res := make([]UWEntry, 0)
	cache := make(map[string]UWEntry)
	for _,v := range entries {
		if _,exist := rp.cache[v.Id]; !exist {
			res = append(res, v)
		}
		cache[v.Id] = v
	}
	rp.cache = cache
	return res
}