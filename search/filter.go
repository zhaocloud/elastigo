// Copyright 2013 Matthew Baird
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package search

import (
	"encoding/json"
	"fmt"
	. "github.com/araddon/gou"
)

var (
	_ = DEBUG
)

// A bool (and/or) clause
type BoolClause string

// Filter clause is either a boolClause or FilterOp
type FilterClause interface {
	String() string
}

// A wrapper to allow for custom serialization
type FilterWrap struct {
	BoolClause string
	Filters    []interface{}
}

func NewFilterWrap() *FilterWrap {
	return &FilterWrap{Filters: make([]interface{}, 0), BoolClause: "and"}
}

func (f *FilterWrap) String() string {
	return fmt.Sprintf(`%s: %d:%v`, f.BoolClause, len(f.Filters), f.Filters)
}

// Custom marshalling to support the query dsl
func (f *FilterWrap) AddFilters(fl ...interface{}) {
	if len(fl) > 1 {
		fc := fl[0]
		switch fc.(type) {
		case BoolClause, string:
			f.BoolClause = fc.(string)
			fl = fl[1:]
		}
	}
	f.Filters = append(f.Filters, fl...)
}

// Custom marshalling to support the query dsl
func (f *FilterWrap) MarshalJSON() ([]byte, error) {
	var root interface{}
	if len(f.Filters) > 1 {
		if f.BoolClause == "bool" {
			filter := map[string]interface{}{}
			for _, flt := range f.Filters {
				filter[flt.(*FilterWrap).BoolClause] = flt.(*FilterWrap).Filters
			}
			root = map[string]interface{}{f.BoolClause: filter}
		} else {
			root = map[string]interface{}{f.BoolClause: f.Filters}
		}

	} else if len(f.Filters) == 1 {
		switch f.BoolClause {
		case "bool", "must", "should", "must_not":
			root = map[string]interface{}{f.BoolClause: f.Filters[0]}
		default:
			root = f.Filters[0]
		}
	}
	return json.Marshal(root)
}

/*
	"filter": {
		"range": {
		  "@timestamp": {
		    "from": "2012-12-29T16:52:48+00:00",
		    "to": "2012-12-29T17:52:48+00:00"
		  }
		}
	}
	"filter": {
	    "missing": {
	        "field": "repository.name"
	    }
	}

	"filter" : {
	    "terms" : {
	        "user" : ["kimchy", "elasticsearch"],
	        "execution" : "bool",
	        "_cache": true
	    }
	}

	"filter" : {
	    "term" : { "user" : "kimchy"}
	}

	"filter" : {
	    "and" : [
	        {
	            "range" : {
	                "postDate" : {
	                    "from" : "2010-03-01",
	                    "to" : "2010-04-01"
	                }
	            }
	        },
	        {
	            "prefix" : { "name.second" : "ba" }
	        }
	    ]
	}

*/

// Filter Operation
//
//   Filter().Term("user","kimchy")
//
//   // we use variadics to allow n arguments, first is the "field" rest are values
//   Filter().Terms("user", "kimchy", "elasticsearch")
//
//   Filter().Exists("repository.name")
//
func Filter() *FilterOp {
	return &FilterOp{}
}

func CompoundFilter(fl ...interface{}) *FilterWrap {
	FilterVal := NewFilterWrap()
	FilterVal.AddFilters(fl...)
	return FilterVal
}

type FilterOp struct {
	curField    string
	TermsMap    map[string][]interface{}          `json:"terms,omitempty"`
	TermMap     map[string]interface{}            `json:"term,omitempty"`
	Range       map[string]map[string]interface{} `json:"range,omitempty"`
	Exist       map[string]string                 `json:"exists,omitempty"`
	MisssingVal map[string]string                 `json:"missing,omitempty"`
}

// A range is a special type of Filter operation
//
//    Range().Exists("repository.name")
func Range() *FilterOp {
	return &FilterOp{Range: make(map[string]map[string]interface{})}
}

func (f *FilterOp) Field(fld string) *FilterOp {
	f.curField = fld
	if _, ok := f.Range[fld]; !ok {
		m := make(map[string]interface{})
		f.Range[fld] = m
	}
	return f
}

// Filter Terms
//
//   Filter().Terms("user","kimchy")
//
//   // we use variadics to allow n arguments, first is the "field" rest are values
//   Filter().Terms("user", "kimchy", "elasticsearch")
//
func (f *FilterOp) Terms(field string, values ...interface{}) *FilterOp {
	if len(f.TermsMap) == 0 {
		f.TermsMap = make(map[string][]interface{})
	}
	for _, val := range values {
		f.TermsMap[field] = append(f.TermsMap[field], val)
	}

	return f
}
func (f *FilterOp) Term(field string, value interface{}) *FilterOp {
	if len(f.TermMap) == 0 {
		f.TermMap = make(map[string]interface{})
	}
	f.TermMap[field] = value
	return f
}
func (f *FilterOp) From(from string) *FilterOp {
	f.Range[f.curField]["from"] = from
	return f
}
func (f *FilterOp) To(to string) *FilterOp {
	f.Range[f.curField]["to"] = to
	return f
}
func (f *FilterOp) Gt(gt int) *FilterOp {
	f.Range[f.curField]["gt"] = float64(gt)
	return f
}
func (f *FilterOp) Exists(name string) *FilterOp {
	f.Exist = map[string]string{"field": name}
	return f
}
func (f *FilterOp) Missing(name string) *FilterOp {
	f.MisssingVal = map[string]string{"field": name}
	return f
}

// Add another Filterop, "combines" two filter ops into one
func (f *FilterOp) Add(fop *FilterOp) *FilterOp {
	// TODO, this is invalid, refactor
	if len(fop.Exist) > 0 {
		f.Exist = fop.Exist
	}
	if len(fop.MisssingVal) > 0 {
		f.MisssingVal = fop.MisssingVal
	}
	if len(fop.Range) > 0 {
		f.Range = fop.Range
	}
	return f
}
