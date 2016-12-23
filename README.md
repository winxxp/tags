#tags

Encode struct field to url parameters by tags

## Quick Start

```Go

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

u := tags.URLValue{TagName: "tag"}
str := u.Encode(&device)

/*
//output: 
channel%5B0%5D.dc=0&channel%5B0%5D.eu=mV&channel%5B0%5D.gain=1&channel%5B0%5D.name=ch1&channel%5B1%5D.dc=0&channel%5B1%5D.eu=mV&channel%5B1%5D.gain=1&channel%5B1%5D.name=ch1&mode=2&name=my-device&remark.build=2016-12-23+14%3A57%3A24.1324366+%2B0800+CST&remark.version=1.0.1&sn=dd-aa-xx

//url.Values: 
channel[0].gain=[1]
channel[0].name=[ch1]
channel[1].dc=[0]
channel[1].gain=[1]
mode=[2]
name=[my-device]
channel[0].dc=[0]
channel[0].eu=[mV]
sn=[dd-aa-xx]
remark.build=[2016-12-23 14:57:24.1324366 +0800 CST]
remark.version=[1.0.1]
channel[1].eu=[mV]
channel[1].name=[ch1]
*/
```

