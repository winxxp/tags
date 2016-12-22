package tags

import "testing"

type S2 struct {
	SF1 int `mts:"sf1"`

	Data []int `mts:"data"`
}

type MockS struct {
	Field1 int            `mts:"field1"`
	Field2 string         `mts:"field2"`
	Field3 S2             `mts:"field3"`
	Field4 map[string]int `mts:"field3"`
}

func TestEncode(t *testing.T) {
	s := MockS{
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
	}

	u := URLValue{TagName: "mts"}
	str := u.Encode(&s)

	t.Log(str)
}
