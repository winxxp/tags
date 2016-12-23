package tags

import "testing"

type S2 struct {
	SF1 int `mts:"sf1"`

	Data []int `mts:"data"`
}

type MockS struct {
	Field1 int            `mts:"int"`
	Field2 string         `mts:"string"`
	Field3 S2             `mts:"struct"`
	Field4 map[string]int `mts:"map"`
	Field5 []S2           `mts:"structPtr"`
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
	}

	u := URLValue{TagName: "mts"}
	str := u.Encode(&s)

	t.Log(str)
}
