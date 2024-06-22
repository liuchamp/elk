package annotation

import (
	"encoding/json"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema"
	"fmt"
)

type Query struct {
	// 默认查询没有前缀， 操作为eq
	Range   []string `json:"range"`   // 支持的查询方式,gt, gte, lt, lte
	Regex   bool     `json:"regex"`   // 当类型为string时， 支付支持正则， 正则前缀为 re_
	Contain bool     `json:"contain"` // 当类型string时， 支持包含操作取代 eq
}

func (q Query) Name() string {
	return "elk_query"
}

// Decode from ent.
func (a *Query) Decode(o interface{}) error {
	buf, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, a)
}

func QueryForOperation(n *gen.Field) *Query {
	// If there are no annotations given do not load any groups.
	ant := &Query{}
	if n.Annotations == nil || n.Annotations[ant.Name()] == nil {
		return nil
	}
	// Decode the types annotation and extract the groups requested for the given operation.
	if err := ant.Decode(n.Annotations[ant.Name()]); err != nil {
		fmt.Errorf("cat not decode data")
		return nil
	}

	return ant
}

// Merge implements ent.Merger interface.
func (a Query) Merge(other schema.Annotation) schema.Annotation {
	var ant Query
	switch o := other.(type) {
	case Query:
		ant = o
	case *Query:
		if o != nil {
			ant = *o
		}
	default:
		return a
	}

	if len(ant.Range) > 0 {
		a.Range = ant.Range
	}
	if ant.Contain {
		a.Contain = ant.Contain
	}
	if ant.Regex {
		a.Regex = ant.Regex
	}

	return a
}

//// MarshalJSON implements the json.Marshaler interface.
//func (f Query) MarshalJSON() ([]byte, error) {
//	type Alias Query
//	return json.Marshal(&struct{ *Alias }{Alias: (*Alias)(&f)})
//}
//
//// UnmarshalJSON implements the json.Unmarshaler interface.
//func (f *Query) UnmarshalJSON(data []byte) error {
//	type Alias Query
//	aux := &struct{ *Alias }{Alias: (*Alias)(f)}
//	return json.Unmarshal(data, aux)
//}

const (
	RangeGt  = "gt"
	RangeGte = "gte"
	RangeLt  = "lt"
	RangeLte = "lte"
)

var _ interface {
	schema.Annotation
	schema.Merger
} = (*Query)(nil)
