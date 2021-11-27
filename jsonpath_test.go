package jsonpath

import (
	"encoding/json"
	"fmt"
	"testing"
)

type obj struct {
	Title  string
	Saunas []sauna
}

type sauna struct {
	Name         string
	Location     string
	Closed       []string
	BathSections []bathSection
}

type bathSection struct {
	Type       string
	SaunaRooms []saunaRoom
}

type saunaRoom struct {
	People      int
	Temperature float32
	Tags        []string
}

var data = obj{
	Title: "タイトル",
	Saunas: []sauna{
		{
			Name:     "草加健康センター",
			Location: "埼玉県",
			Closed:   []string{},
			BathSections: []bathSection{
				{
					Type: "男",
					SaunaRooms: []saunaRoom{
						{
							People:      28,
							Temperature: 92.0,
							Tags:        []string{"ドライ", "ロウリュ"},
						},
					},
				},
				{
					Type: "女",
					SaunaRooms: []saunaRoom{
						{
							People:      25,
							Temperature: 80.0,
							Tags:        []string{"ドライ", "ロウリュ"},
						},
						{
							People:      4,
							Temperature: 54.0,
							Tags:        []string{"スチーム"},
						},
					},
				},
			},
		},
		{
			Name:     "金春湯",
			Location: "東京都",
			Closed:   []string{"火曜日"},
			BathSections: []bathSection{
				{
					Type: "男",
					SaunaRooms: []saunaRoom{
						{
							People:      6,
							Temperature: 90.0,
							Tags:        []string{"ドライ"},
						},
					},
				},
				{
					Type: "女",
					SaunaRooms: []saunaRoom{
						{
							People:      6,
							Temperature: 92.0,
							Tags:        []string{"ドライ"},
						},
					},
				},
			},
		},
	},
}

func toInterfaceJSON(input interface{}) (interface{}, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	var d interface{}
	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func TestGet(t *testing.T) {
	query := "$.Title"
	p := NewPath(query)

	err := p.Parse()
	if err != nil {
		t.Error(err)
	}

	input, err := toInterfaceJSON(data)
	if err != nil {
		t.Error(err)
	}

	output, err := p.Execute(input)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%#v\n", output)

	if output_v, ok := output.(string); !ok || output_v != "タイトル" {
		t.Errorf("TestGet error: %v is not match タイトル", output_v)
	}

	return
}
