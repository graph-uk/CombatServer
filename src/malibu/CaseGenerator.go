package main

import (
	"log"
)

type enumerableElement struct {
	name                string
	variants            []string
	currentVariantIndex int
}

func (t *enumerableElement) LoadVariants(variants []string) {
	if len(variants) < 1 {
		log.Fatal("Enum element should contain at least 1 variant")
	}
	t.variants = variants
}

// select next element. Return true if overflow and start from begin
func (t *enumerableElement) Next() bool {
	if t.currentVariantIndex == len(t.variants)-1 {
		t.currentVariantIndex = 0
		return true
	} else {
		t.currentVariantIndex++
		return false
	}
}

type combinator struct {
	elements []enumerableElement
}

func (t *combinator) Next() (map[string]string, bool) {
	var result map[string]string
	result = make(map[string]string)

	for _, curElement := range t.elements {
		result[curElement.name] = curElement.variants[curElement.currentVariantIndex]
	}

	allEnd := false
	for curElementIndex, _ := range t.elements {
		if !t.elements[curElementIndex].Next() {
			break
		} else {
			if curElementIndex == len(t.elements)-1 {
				allEnd = true
			}
		}
	}

	return result, allEnd
}

func (t *combinator) LoadElements(input map[string][]string) {
	for curElementName, curElement := range input {
		var element enumerableElement
		element.LoadVariants(curElement)
		element.name = curElementName
		t.elements = append(t.elements, element)
	}
}

func getAllParamsCombinations(input map[string][]string) []*map[string]string {
	var result []*map[string]string

	var combinator combinator
	combinator.LoadElements(input)

	for {
		curCombination, isEnd := combinator.Next()
		result = append(result, &curCombination)

		if isEnd {
			break
		}
	}
	return result
}
