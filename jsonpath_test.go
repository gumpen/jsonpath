package jsonpath

import (
	"encoding/json"
	"fmt"
	"testing"
)

// type obj struct {
// 	Title  string
// 	Saunas []sauna
// }

// type sauna struct {
// 	Name         string
// 	Location     string
// 	Closed       []string
// 	BathSections []bathSection
// }

// type bathSection struct {
// 	Type       string
// 	SaunaRooms []saunaRoom
// }

// type saunaRoom struct {
// 	People      int
// 	Temperature float32
// 	Tags        []string
// }

// var data = obj{
// 	Title: "タイトル",
// 	Saunas: []sauna{
// 		{
// 			Name:     "草加健康センター",
// 			Location: "埼玉県",
// 			Closed:   []string{},
// 			BathSections: []bathSection{
// 				{
// 					Type: "男",
// 					SaunaRooms: []saunaRoom{
// 						{
// 							People:      28,
// 							Temperature: 92.0,
// 							Tags:        []string{"ドライ", "ロウリュ"},
// 						},
// 					},
// 				},
// 				{
// 					Type: "女",
// 					SaunaRooms: []saunaRoom{
// 						{
// 							People:      25,
// 							Temperature: 80.0,
// 							Tags:        []string{"ドライ", "ロウリュ"},
// 						},
// 						{
// 							People:      4,
// 							Temperature: 54.0,
// 							Tags:        []string{"スチーム"},
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			Name:     "金春湯",
// 			Location: "東京都",
// 			Closed:   []string{"火曜日"},
// 			BathSections: []bathSection{
// 				{
// 					Type: "男",
// 					SaunaRooms: []saunaRoom{
// 						{
// 							People:      6,
// 							Temperature: 90.0,
// 							Tags:        []string{"ドライ"},
// 						},
// 					},
// 				},
// 				{
// 					Type: "女",
// 					SaunaRooms: []saunaRoom{
// 						{
// 							People:      6,
// 							Temperature: 92.0,
// 							Tags:        []string{"ドライ"},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},
// }

var data = `{
    "Title": "タイトル",
    "Author": {
        "Name": "gumpen",
        "ID": 123456
    },
    "Saunas": [
        {
            "Name": "草加健康センター",
            "Location": "埼玉県",
            "Closed": [],
            "BathSections": [
                {
                    "Type": "男",
                    "SaunaRooms": [
                        {
                            "People": 28,
                            "Temperature": 92,
                            "Tags": [
                                "ドライ",
                                "ロウリュ"
                            ]
                        }
                    ]
                },
                {
                    "Type": "女",
                    "SaunaRooms": [
                        {
                            "People": 25,
                            "Temperature": 80,
                            "Tags": [
                                "ドライ",
                                "ロウリュ"
                            ]
                        },
                        {
                            "People": 4,
                            "Temperature": 54,
                            "Tags": [
                                "スチーム"
                            ]
                        }
                    ]
                }
            ]
        },
        {
            "Name": "金春湯",
            "Location": "東京都",
            "Closed": [
                "火曜日"
            ],
            "BathSections": [
                {
                    "Type": "男",
                    "SaunaRooms": [
                        {
                            "People": 6,
                            "Temperature": 90,
                            "Tags": [
                                "ドライ"
                            ]
                        }
                    ]
                },
                {
                    "Type": "女",
                    "SaunaRooms": [
                        {
                            "People": 6,
                            "Temperature": 92,
                            "Tags": [
                                "ドライ"
                            ]
                        }
                    ]
                }
            ]
        }
    ]
}`

// func toInterfaceJSON(input interface{}) (interface{}, error) {
// 	b, err := json.Marshal(input)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var d interface{}
// 	err = json.Unmarshal(b, &d)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return d, nil
// }

type testCase struct {
	name   string
	query  string
	input  interface{}
	expect interface{}
}

func testJSONPath(cases []testCase, t *testing.T) {
	for _, c := range cases {
		defer func() {
			err := recover()
			if err != nil {
				t.Errorf("test %s error: Recover!:%s", c.name, err)
			}
		}()

		p := NewPath(c.query)
		err := p.Parse()
		if err != nil {
			t.Error(err)
		}

		output, err := p.Execute(c.input)
		if err != nil {
			t.Error(err)
		}

		sOutput := fmt.Sprint(output)
		sExpect := fmt.Sprint(c.expect)

		if sOutput != sExpect {
			t.Errorf("test %s error: output %s should be expected %s", c.name, output, c.expect)
		}

	}

	return
}

func TestGet(t *testing.T) {
	var d interface{}
	err := json.Unmarshal([]byte(data), &d)
	if err != nil {
		t.Error(err)
	}

	cases := []testCase{
		{
			name:   "child element",
			query:  "$.Title",
			input:  d,
			expect: "タイトル",
		},
		{
			name:   "child element 2",
			query:  "$.Author.ID",
			input:  d,
			expect: float64(123456), // goでJSONのnumberはfloat64型になる
		},
		{
			name:   "index access",
			query:  "$.Saunas[0].Name",
			input:  d,
			expect: "草加健康センター",
		},
		{
			name:   "union",
			query:  "$.Saunas[0,1].Name",
			input:  d,
			expect: []string{"草加健康センター", "金春湯"},
		},
		{
			name:   "multi union",
			query:  "$.Saunas[0,1].BathSections[0,1].Type",
			input:  d,
			expect: []string{"男", "女", "男", "女"},
		},
		{
			name:   "slice",
			query:  "$.Saunas[0:1].Location",
			input:  d,
			expect: "埼玉県",
		},
		{
			name:   "slice 2",
			query:  "$.Saunas[0:2].Location",
			input:  d,
			expect: []string{"埼玉県", "東京都"},
		},
		{
			name:   "slice blank start",
			query:  "$.Saunas[:2].Location",
			input:  d,
			expect: []string{"埼玉県", "東京都"},
		},
		{
			name:   "slice blank end",
			query:  "$.Saunas[1:].Location",
			input:  d,
			expect: "東京都",
		},
		{
			name:   "slice blank both",
			query:  "$.Saunas[:].Location",
			input:  d,
			expect: []string{"埼玉県", "東京都"},
		},
		{
			name:   "union and slice",
			query:  "$.Saunas[0,1].BathSections[:].Type",
			input:  d,
			expect: []string{"男", "女", "男", "女"},
		},
	}

	testJSONPath(cases, t)

	return
}
