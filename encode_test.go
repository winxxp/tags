package tags

import (
	"net/url"
	"reflect"
	"testing"
)

type Device struct {
	SN      string            `tag:"sn,name=DeviceSN" json:"sn"`
	Name    string            `tag:"name,name=DeviceName" json:"name"`
	Mode    int               `tag:"mode,name=DeviceMode" json:"mode"`
	Channel []Channel         `tag:"channel,name=DeviceChannel" json:"channel"`
	Remark  map[string]string `tag:"remark,name=DeviceRemark" json:"remark"`
}

type Channel struct {
	Name string  `tag:"name,name=ChannelName" json:"name"`
	Eu   string  `tag:"eu" json:"eu"`
	DC   float32 `tag:"dc" json:"dc"`
	Gain float32 `tag:"gain" json:"gain"`
}

var device = Device{
	SN:   "dd-aa-xx",
	Name: "my-device",
	Mode: 2,
	Channel: []Channel{
		{
			Name: "ch1",
			Eu:   "mV",
			DC:   0,
			Gain: 1,
		},
		{
			Name: "ch1",
			Eu:   "mV",
			DC:   0,
			Gain: 1,
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

	exptected, err := url.ParseQuery(str)
	if err != nil {
		t.Error(err)
		return
	}

	var result = url.Values{
		"channel[0].eu":   []string{"mV"},
		"channel[1].dc":   []string{"0"},
		"channel[1].eu":   []string{"mV"},
		"channel[1].name": []string{"ch1"},
		"name":            []string{"my-device"},
		"channel[0].dc":   []string{"0"},
		"channel[0].gain": []string{"1"},
		"channel[0].name": []string{"ch1"},
		"channel[1].gain": []string{"1"},
		"mode":            []string{"2"},
		"remark.version":  []string{"1.0.1"},
	}

	if !reflect.DeepEqual(result, exptected) {
		t.Fail()
	}

	for k, v := range exptected {
		t.Logf("%s=%s\n", k, v)
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
	u := New("tag", &SubTag{"name"})
	v := u.Values(&device)
	str := v.Encode()

	t.Log(v)
	t.Log(str)

	values, err := url.ParseQuery(str)
	if err != nil {
		t.Error(err)
		return
	}

	expected := url.Values{
		"DeviceChannel[1].dc":          []string{"0"},
		"DeviceChannel[1].eu":          []string{"mV"},
		"DeviceName":                   []string{"my-device"},
		"DeviceRemark.version":         []string{"1.0.1"},
		"DeviceChannel[0].ChannelName": []string{"ch1"},
		"DeviceChannel[0].dc":          []string{"0"},
		"DeviceChannel[0].eu":          []string{"mV"},
		"DeviceChannel[1].ChannelName": []string{"ch1"},
		"DeviceChannel[0].gain":        []string{"1"},
		"DeviceChannel[1].gain":        []string{"1"},
		"DeviceMode":                   []string{"2"},
		"DeviceSN":                     []string{"dd-aa-xx"},
	}

	if !reflect.DeepEqual(values, expected) {
		t.Fail()
	}

	for k, v := range values {
		t.Logf("%s=%s\n", k, v)
	}
}
