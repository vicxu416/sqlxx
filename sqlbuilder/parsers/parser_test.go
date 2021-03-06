package parsers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type source struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"time"`
	Ignore    string
}

func TestParse(t *testing.T) {
	var testcases = []struct {
		name       string
		source     interface{}
		contains   []string
		len        int
		allowEmpty bool
	}{
		{name: "struct_1", source: source{ID: 1, Name: "test"}, contains: []string{"id", "name"}, len: 2, allowEmpty: false},
		{name: "struct_2", source: source{ID: 1, Name: "test", CreatedAt: time.Now()}, contains: []string{"id", "name", "time"}, len: 3, allowEmpty: false},
		{name: "struct_3", source: source{ID: 0, Name: "", CreatedAt: time.Now()}, contains: []string{"id", "name", "time"}, len: 3, allowEmpty: true},
		{name: "slice_struct_1", source: []source{{ID: 1, Name: "test", CreatedAt: time.Now()}}, contains: []string{"id", "name", "time"}, len: 3, allowEmpty: false},
		{name: "slice_struct_2", source: []source{{ID: 1, Name: "test", CreatedAt: time.Now()}, {ID: 1, Name: "test", CreatedAt: time.Now()}}, contains: []string{"id", "name", "time"}, len: 6, allowEmpty: false},
		{name: "map_1", source: map[string]interface{}{"id": 1, "name": ""}, contains: []string{"id", "name"}, len: 2, allowEmpty: true},
		{name: "map_2", source: []map[string]interface{}{{"id": 1, "name": ""}}, contains: []string{"id", "name"}, len: 2, allowEmpty: true},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			parser, err := New(testcase.source, testcase.allowEmpty)
			assert.Nil(t, err)
			for _, contain := range testcase.contains {
				assert.Contains(t, parser.Fields, contain)
				assert.Contains(t, parser.NamedValues[0], ":"+contain)
			}
			assert.Len(t, parser.Data, testcase.len)
		})
	}
}

func TestQueryParser(t *testing.T) {
	var testcases = []struct {
		name   string
		source interface{}
		expect map[string]interface{}
	}{
		{name: "struct_1", source: source{ID: 1, Name: "test"}, expect: map[string]interface{}{"id": int64(1), "name": "test"}},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			conditon := make(map[string]interface{})
			err := ParseConditions(testcase.source, conditon, false)
			assert.Nil(t, err)
			assert.Equal(t, testcase.expect, conditon)
		})
	}
}
