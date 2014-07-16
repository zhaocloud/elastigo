package search

func HighlightWith(fields ...string) *Highlighting {
	h := &Highlighting{Fields: map[string]interface{}{}}
	for _, f := range fields {
		h.Fields[f] = make(map[string]interface{})
	}
	return h
}

type Highlighting struct {
	PreTags  []string               `json:"pre_tags,omitempty"`
	PostTags []string               `json:"post_tags,omitempty"`
	Fields   map[string]interface{} `json:"fields"`
}

func (h *Highlighting) AddPreTags(tag ...string) *Highlighting {
	h.PreTags = append(h.PreTags, tag...)
	return h
}

func (h *Highlighting) AddPostTags(tag ...string) *Highlighting {
	h.PostTags = append(h.PostTags, tag...)
	return h
}
