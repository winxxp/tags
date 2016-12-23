package tags

import (
	"net/url"
	"testing"
)

type S2 struct {
	SF1  int   `tag:"sf1"`
	Data []int `tag:"data"`
}

type MockS struct {
	Field1 int            `tag:"int"`
	Field2 string         `tag:"string"`
	Field3 S2             `tag:"struct"`
	Field4 map[string]int `tag:"map"`
	Field5 []S2           `tag:"structArray"`
	Field6 *S2            `tag:"structPtr"`
	Field7 *S2            `tag:"structPtr"`
}

var data = MockS{
	Field1: 1,
	Field2: "test",
	Field3: S2{
		SF1:  2,
		Data: []int{10, 11, 12},
	},
	Field4: map[string]int{
		"map1": 100,
		"map0": 99,
	},
	Field5: []S2{
		S2{
			SF1:  2,
			Data: []int{20, 21},
		},
		S2{
			SF1:  3,
			Data: []int{30, 31},
		},
	},
	Field6: &S2{
		SF1:  4,
		Data: []int{41, 42},
	},
}

func TestEncode(t *testing.T) {
	u := URLValue{TagName: "tag"}
	str := u.Encode(&data)

	values, err := url.ParseQuery(str)
	if err != nil {
		t.Error(err)
		return
	}

	for k, v := range values {
		t.Logf("%s=%s\n", k, v)
	}
}

func BenchmarkURLValue_Encode(b *testing.B) {
	u := URLValue{TagName: "tag"}
	for i := 0; i < b.N; i++ {
		_ = u.Encode(data)
	}
}

func TestURLValue_Encode1(t *testing.T) {
	u := URLValue{TagName: "tag"}
	a := 1
	str := u.Encode(a)
	t.Log(str)
}
