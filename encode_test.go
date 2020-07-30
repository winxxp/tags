package tags

import (
	"net/url"
	"reflect"
	"testing"
)

type Device struct {
	SN        string            `tag:"sn,name=DeviceSN" json:"sn"`
	Name      string            `tag:"name,name=DeviceName" json:"name"`
	Mode      int               `tag:"mode,name=DeviceMode" json:"mode"`
	StartTime []int             `tag:"startTime,name=StartTime" json:"startTime"`
	Channel   []Channel         `tag:"channel,name=DeviceChannel" json:"channel"`
	Remark    map[string]string `tag:"remark,name=DeviceRemark" json:"remark"`
}

type Channel struct {
	Name string  `tag:"name,name=ChannelName" json:"name"`
	Eu   string  `tag:"eu" json:"eu"`
	DC   float32 `tag:"dc" json:"dc"`
	Gain float32 `tag:"gain" json:"gain"`
}

var device = Device{
	SN:        "dd-aa-xx",
	Name:      "my-device",
	Mode:      2,
	StartTime: []int{1, 2, 3},
	Channel: []Channel{
		{
			Name: "ch1",
			Eu:   "mV",
			DC:   0,
			Gain: 3120372.5,
		},
		{
			Name: "ch2",
			Eu:   "mV",
			DC:   0,
			Gain: 11.1234,
		},
	},
	Remark: map[string]string{
		"version": "1.0.1",
	},
}

func TestEncode(t *testing.T) {
	u := New("json", nil)
	v := u.Values(&device)
	str := v.Encode()

	t.Log(v)
	t.Log(str)

	result, err := url.ParseQuery(str)
	if err != nil {
		t.Error(err)
		return
	}

	for k, v := range result {
		t.Logf("%s=%s\n", k, v)
	}

	var expected = url.Values{
		"name":            []string{"my-device"},
		"mode":            []string{"2"},
		"DeviceSN":        []string{"dd-aa-xx"},
		"channel[0].name": []string{"ch1"},
		"channel[0].eu":   []string{"mV"},
		"channel[0].dc":   []string{"0"},
		"channel[0].gain": []string{"3120372.5"},
		"channel[1].name": []string{"ch2"},
		"channel[1].eu":   []string{"mV"},
		"channel[1].dc":   []string{"0"},
		"channel[1].gain": []string{"11.12339973449707"},
		"remark.version":  []string{"1.0.1"},
		"startTime[0]":    []string{"1"},
		"startTime[1]":    []string{"2"},
		"startTime[2]":    []string{"3"},
	}

	if reflect.DeepEqual(result, expected) {
		t.Fail()
	}

}

func BenchmarkURLValue_Encode(b *testing.B) {
	u := New("json", nil)
	for i := 0; i < b.N; i++ {
		_ = u.Encode(device)
	}
	b.ReportAllocs()
}

func TestURLValue_Encode1(t *testing.T) {
	u := New("tag", NewSubTagFinder("name"))
	v := u.Values(&device)
	str := v.Encode()

	t.Log(v)
	t.Log(str)

	result, err := url.ParseQuery(str)
	if err != nil {
		t.Error(err)
		return
	}
	for k, v := range result {
		t.Logf("%s=%s\n", k, v)
	}

	var expected = url.Values{
		"name":            []string{"my-device"},
		"mode":            []string{"2"},
		"DeviceSN":        []string{"dd-aa-xx"},
		"channel[0].name": []string{"ch1"},
		"channel[0].eu":   []string{"mV"},
		"channel[0].dc":   []string{"0"},
		"channel[0].gain": []string{"3120372.5"},
		"channel[1].name": []string{"ch2"},
		"channel[1].eu":   []string{"mV"},
		"channel[1].dc":   []string{"0"},
		"channel[1].gain": []string{"11.12339973449707"},
		"remark.version":  []string{"1.0.1"},
		"StartTime[0]":    []string{"1"},
		"StartTime[1]":    []string{"2"},
		"StartTime[2]":    []string{"3"},
	}

	if reflect.DeepEqual(result, expected) {
		t.Fail()
	}
}
