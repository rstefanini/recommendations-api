package model

import "sort"

type ranking []ProductSKU
type ProductHits map[ProductSKU]Hits

func (p *ProductHits) GetTop(n int) ranking {
	h := make([]Hits, 0, len(*p))
	for _, v := range *p {
		h = append(h, v)
	}

	sort.Slice(h, func(i, j int) bool {
		return h[i].Hits > h[j].Hits
	})

	r := ranking{}

	for k, v := range h {
		if k == n {
			break
		}

		r = append(r, v.Product)
	}

	return r
}

type Hits struct {
	Product ProductSKU
	Hits    int
}
