package fuzzy

import (
	"unicode/utf8"

	"golang.org/x/text/transform"
	"me.daily/src/transformer"
)

// 模糊搜尋
//	https://github.com/lithammer/fuzzysearch
func FuzzySearch(source, target string, t transform.Transformer) bool {
	if t == nil {
		t = transformer.Nop
	}

	source = stringTransform(source, t)
	target = stringTransform(target, t)

	lenSource := utf8.RuneCountInString(source)
	lenTarget := utf8.RuneCountInString(target)

	// 長度
	if lenSource-lenTarget < 0 {
		return false
	}

	// 來源相同
	if source == target {
		return true
	}

	index := 0
	for _, c1 := range target {

		flag := false
		subStr := source[index:]
		for _, c2 := range subStr {
			index++

			if c1 == c2 {
				flag = true
				break
			}
		}

		if !flag {
			return false
		}
	}

	return true
}

func stringTransform(s string, t transform.Transformer) (str string) {
	var err error
	str, _, err = transform.String(t, s)
	if err != nil {
		str = s
	}
	return
}
