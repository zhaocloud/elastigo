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
	// "fmt"
	//"log"
	"strings"
)

// QueryDsl creates a new Query Dsl
func Query() *QueryDsl {
	return &QueryDsl{}
}

/*

some ways to serialize
"query": {
	"filtered": {
	  "query": {
	    "query_string": {
	      "default_operator": "OR",
	      "default_field": "_all",
	      "query": " actor:\"bob\"  AND type:\"EventType\""
	    }
	  },
	  "filter": {
	    "range": {
	      "@timestamp": {
	        "from": "2012-12-29T16:52:48+00:00",
	        "to": "2012-12-29T17:52:48+00:00"
	      }
	    }
	  }
	}
},

"query" : {
    "term" : { "user" : "kimchy" }
}

"query" : {
    "match_all" : {}
},
*/
type QueryDsl struct {
	QueryEmbed
	QueryBody map[string]interface{}
	// FilterVal *FilterOp `json:"filter,omitempty"`
}

// The core Query Syntax can be embedded as a child of a variety of different parents
type QueryEmbed struct {
	MatchAll *MatchAll         `json:"match_all,omitempty"`
	Terms    map[string]string `json:"term,omitempty"`
	Qs       *QueryString      `json:"query_string,omitempty"`
	//Exist    string            `json:"_exists_,omitempty"`
}

func (qd *QueryDsl) HasQuery() bool {
	q := qd.QueryEmbed
	if q.Qs != nil || len(q.Terms) > 0 || q.MatchAll != nil {
		return true
	}
	return false
}

// MarshalJSON provides custom marshalling to support the query dsl which is a conditional
// json format, not always the same parent/children
func (qd *QueryDsl) MarshalJSON() ([]byte, error) {
	if qd.QueryBody != nil {
		return json.Marshal(qd.QueryBody)
	}
	return json.Marshal(qd.QueryEmbed)
}

// get all
func (q *QueryDsl) All() *QueryDsl {
	q.MatchAll = &MatchAll{""}
	return q
}

// // Limit the query to this range
// func (q *QueryDsl) Range(fop *FilterOp) *QueryDsl {
// 	if q.FilterVal == nil {
// 		q.FilterVal = fop
// 		return q
// 	}
// 	// TODO:  this is not valid, refactor
// 	q.FilterVal.Add(fop)
// 	return q
// }

// Add a term search for a specific field
//    Term("user","kimchy")
func (q *QueryDsl) Term(name, value string) *QueryDsl {
	if len(q.Terms) == 0 {
		q.Terms = make(map[string]string)
	}
	q.Terms[name] = value
	return q
}

// The raw search strings (lucene valid)
func (q *QueryDsl) Search(searchFor string) *QueryDsl {
	//I don't think this is right, it is not a filter.query, it should be q query?
	qs := NewQueryString("", "")
	q.QueryEmbed.Qs = &qs
	q.QueryEmbed.Qs.Query = searchFor
	return q
}

// Querystring operations
func (q *QueryDsl) Qs(qs *QueryString) *QueryDsl {
	q.QueryEmbed.Qs = qs
	return q
}

// Fields in query_string search
//     Fields("fieldname","search_for","","")
//
//     Fields("fieldname,field2,field3","search_for","","")
//
//     Fields("fieldname,field2,field3","search_for","field_exists","")
func (q *QueryDsl) Fields(fields, search, exists, missing string) *QueryDsl {
	fieldList := strings.Split(fields, ",")
	qs := NewQueryString("", "")
	q.QueryEmbed.Qs = &qs
	q.QueryEmbed.Qs.Query = search
	if len(fieldList) == 1 {
		q.QueryEmbed.Qs.DefaultField = fields
	} else {
		q.QueryEmbed.Qs.Fields = fieldList
	}
	q.QueryEmbed.Qs.Exists = exists
	q.QueryEmbed.Qs.Missing = missing
	return q
}

// Filter this query
// func (q *QueryDsl) Filter(f *FilterOp) *QueryDsl {
// 	q.FilterVal = f
// 	return q
// }

type MatchAll struct {
	All string `json:"-"`
}

// should we reuse QueryDsl here?
type QueryWrap struct {
	Qs QueryString `json:"query_string,omitempty"`
}

// QueryString based search
func NewQueryString(field, query string) QueryString {
	return QueryString{"", field, true, query, "", "", nil}
}

type QueryString struct {
	DefaultOperator string   `json:"default_operator,omitempty"`
	DefaultField    string   `json:"default_field,omitempty"`
	AnalyzeWildcard bool     `json:"analyze_wildcard"`
	Query           string   `json:"query,omitempty"`
	Exists          string   `json:"_exists_,omitempty"`
	Missing         string   `json:"_missing_,omitempty"`
	Fields          []string `json:"fields,omitempty"`
	//_exists_:field1,
	//_missing_:field1,
}

// Generic Term based (used in query, facet, filter)
type Term struct {
	Terms     Terms       `json:"terms,omitempty"`
	FilterVal *FilterWrap `json:"facet_filter,omitempty"`
}

type Terms struct {
	Fields []string `json:"field,omitempty"`
	Size   string   `json:"size,omitempty"`
	Regex  string   `json:"regex,omitempty"`
}

func NewTerm(fields ...string) *Term {
	m := &Term{Terms{Fields: fields}, nil}
	return m
}

func (s *Term) Filter(fl ...interface{}) *Term {
	if s.FilterVal == nil {
		s.FilterVal = NewFilterWrap()
	}

	s.FilterVal.AddFilters(fl)
	return s
}

// Custom marshalling
func (t *Terms) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	// TODO:  this isn't getting called!?
	if len(t.Fields) == 1 {
		m["field"] = t.Fields[0]
	} else if len(t.Fields) > 1 {
		m["fields"] = t.Fields
	}
	if len(t.Regex) > 0 {
		m["regex"] = t.Regex
	}
	if len(t.Size) > 0 {
		m["size"] = t.Size
	}
	return json.Marshal(m)
}
