package tags

import (
	"net/url"
	"testing"
	"time"
)

type Device struct {
	SN      string            `tag:"sn"`
	Name    string            `tag:"name"`
	Mode    int               `tag:"mode"`
	Channel []Channel         `tag:"channel"`
	Remark  map[string]string `tag:"remark"`
}

type Channel struct {
	Name string  `tag:"name"`
	Eu   string  `tag:"eu"`
	DC   float32 `tag:"dc"`
	Gain float32 `tag:"gain"`
}

var device = Device{
	SN:   "dd-aa-xx",
	Name: "my-device",
	Mode: 2,
	Channel: []Channel{
		Channel{
			Name: "ch1",
			Eu:   "mV",
			DC:   0,
			Gain: 1,
		},
		Channel{
			Name: "ch1",
			Eu:   "mV",
			DC:   0,
			Gain: 1,
		},
	},
	Remark: map[string]string{
		"build":   time.Now().String(),
		"version": "1.0.1",
	},
}

func TestEncode(t *testing.T) {
	u := Enc{TagName: "tag"}
	v := u.Values(&device)
	str := v.Encode()

	t.Log(v)
	t.Log(str)

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
	u := Enc{TagName: "tag"}
	for i := 0; i < b.N; i++ {
		_ = u.Encode(device)
	}
}

func TestURLValue_Encode1(t *testing.T) {
	u := Enc{TagName: "tag"}
	a := 1
	str := u.Encode(a)
	t.Log(str)
}
