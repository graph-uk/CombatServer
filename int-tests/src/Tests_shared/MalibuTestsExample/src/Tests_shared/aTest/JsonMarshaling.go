package aTest

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type param struct {
	Name     string
	Type     string
	Variants []string
}

type JsonParamOutput struct {
	Params []param
	Tags   []string
}

func (a *ATest) PrintParamsAndTagsInJSON(params interface{}, tags []string) {
	//	m := JsonParamOutput{
	//		Params: []param{
	//			param{
	//				Name:     "SessionTimestamp",
	//				Type:     "StringParam",
	//				Variants: nil,
	//			},
	//			param{
	//				Name:     "Locales",
	//				Type:     "EnumParam",
	//				Variants: []string{"EN", "RU"},
	//			},
	//			param{
	//				Name:     "AdminName",
	//				Type:     "StringParam",
	//				Variants: nil,
	//			},
	//		},
	//		Tags: []string{"lyn2x"},
	//	}
	var m JsonParamOutput

	for _, curTag := range a.Tags {
		m.Tags = append(m.Tags, curTag)
	}

	s := reflect.ValueOf(params).Elem()
	for curFieldIndex := 0; curFieldIndex < s.NumField(); curFieldIndex++ {
		curParam := param{}
		curField := s.Field(curFieldIndex)
		curParam.Name = s.Type().Field(curFieldIndex).Name
		if curField.Type().Name() == "StringParam" {
			curParam.Type = "StringParam"
		}
		if curField.Type().Name() == "EnumParam" {
			curParam.Type = "EnumParam"
			for EnumIndex := 0; EnumIndex < curField.Field(1).Len(); EnumIndex++ {
				curParam.Variants = append(curParam.Variants, curField.Field(1).Index(EnumIndex).String())
			}
		}
		m.Params = append(m.Params, curParam)
	}

	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	var m2 JsonParamOutput
	err = json.Unmarshal(b, &m2)
}
