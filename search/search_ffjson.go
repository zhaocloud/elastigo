// DO NOT EDIT!
// Code generated by ffjson <https://github.com/pquerna/ffjson>
// source: search.go
// DO NOT EDIT!

package search

import (
	"bytes"
	"encoding/json"
)

func (mj *SearchDsl) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(1024)
	err := mj.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (mj *SearchDsl) MarshalJSONBuf(buf *bytes.Buffer) error {
	var err error
	var obj []byte
	var first bool = true
	_ = obj
	_ = err
	_ = first
	buf.WriteString(`{`)
	if len(mj.AggregatesVal) != 0 {
		if first == true {
			first = false
		} else {
			buf.WriteString(`,`)
		}
		buf.WriteString(`"aggregations":`)
		/* Falling back. type=map[string]*search.AggregateDsl kind=map */
		obj, err = json.Marshal(mj.AggregatesVal)
		if err != nil {
			return err
		}
		buf.Write(obj)
	}
	if mj.FacetVal != nil {
		if first == true {
			first = false
		} else {
			buf.WriteString(`,`)
		}
		buf.WriteString(`"facets":`)
		obj, err = mj.FacetVal.MarshalJSON()
		if err != nil {
			return err
		}
		buf.Write(obj)
	}
	if mj.QueryVal != nil {
		if first == true {
			first = false
		} else {
			buf.WriteString(`,`)
		}
		buf.WriteString(`"query": {"filtered":{"query":`)
		obj, err = mj.QueryVal.MarshalJSON()
		if err != nil {
			return err
		}
		buf.Write(obj)
		if mj.FilterVal != nil {
			if first == true {
				first = false
			} else {
				buf.WriteString(`,`)
			}
			buf.WriteString(`"filter":`)
			obj, err = mj.FilterVal.MarshalJSON()
			if err != nil {
				return err
			}
			buf.Write(obj)
		}
		buf.WriteString(`}}`)
	}
	if mj.FromVal != 0 {
		if first == true {
			first = false
		} else {
			buf.WriteString(`,`)
		}
		buf.WriteString(`"from":`)
		ffjson_FormatBits(buf, uint64(mj.FromVal), 10, mj.FromVal < 0)
	}
	if mj.SizeVal != 0 {
		if first == true {
			first = false
		} else {
			buf.WriteString(`,`)
		}
		buf.WriteString(`"size":`)
		ffjson_FormatBits(buf, uint64(mj.SizeVal), 10, mj.SizeVal < 0)
	}
	if len(mj.SortBody) != 0 {
		if first == true {
			first = false
		} else {
			buf.WriteString(`,`)
		}
		buf.WriteString(`"sort":`)
		if mj.SortBody != nil {
			buf.WriteString(`[`)
			for i, v := range mj.SortBody {
				if i != 0 {
					buf.WriteString(`,`)
				}
				obj, err = v.MarshalJSON()
				if err != nil {
					return err
				}
				buf.Write(obj)
			}
			buf.WriteString(`]`)
		} else {
			buf.WriteString(`null`)
		}
	}
	if mj.HighlightVal != nil {
		if first == true {
			first = false
		} else {
			buf.WriteString(`,`)
		}
		buf.WriteString(`"highlight":`)
		obj, err = json.Marshal(mj.HighlightVal)
		if err != nil {
			return err
		}
		buf.Write(obj)
	}
	buf.WriteString(`}`)
	return nil
}

func ffjson_FormatBits(dst *bytes.Buffer, u uint64, base int, neg bool) {
	const (
		digits   = "0123456789abcdefghijklmnopqrstuvwxyz"
		digits01 = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
		digits10 = "0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999"
	)

	var shifts = [len(digits) + 1]uint{
		1 << 1: 1,
		1 << 2: 2,
		1 << 3: 3,
		1 << 4: 4,
		1 << 5: 5,
	}

	if base < 2 || base > len(digits) {
		panic("strconv: illegal AppendInt/FormatInt base")
	}

	var a [64 + 1]byte
	i := len(a)

	if neg {
		u = -u
	}

	if base == 10 {

		for u >= 100 {
			i -= 2
			q := u / 100
			j := uintptr(u - q*100)
			a[i+1] = digits01[j]
			a[i+0] = digits10[j]
			u = q
		}
		if u >= 10 {
			i--
			q := u / 10
			a[i] = digits[uintptr(u-q*10)]
			u = q
		}

	} else if s := shifts[base]; s > 0 {

		b := uint64(base)
		m := uintptr(b) - 1
		for u >= b {
			i--
			a[i] = digits[uintptr(u)&m]
			u >>= s
		}

	} else {

		b := uint64(base)
		for u >= b {
			i--
			a[i] = digits[uintptr(u%b)]
			u /= b
		}
	}

	i--
	a[i] = digits[uintptr(u)]

	if neg {
		i--
		a[i] = '-'
	}

	dst.Write(a[i:])

	return
}
