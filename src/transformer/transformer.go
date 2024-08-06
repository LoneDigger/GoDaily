package transformer

import (
	"bytes"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// https://github.com/lithammer/fuzzysearch

// 不轉換
var Nop = transform.Nop

// 忽略大小寫
var Fold = unicodeFoldTransformer{}

// Unicode 等價性
//
//	https://learnku.com/docs/go-blog/normalization/6554
var Normalized = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

// 忽略大小寫 + Unicode等價性
var NormalizedFold = transform.Chain(Fold, Normalized)

// 轉成小寫
type unicodeFoldTransformer struct{}

func (unicodeFoldTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	runes := bytes.Runes(src)

	lowerRunes := make([]rune, len(runes))
	for i, r := range runes {
		lowerRunes[i] = unicode.ToLower(r)
	}

	srcBytes := []byte(string(lowerRunes))
	n := copy(dst, srcBytes)
	if n < len(srcBytes) {
		err = transform.ErrShortDst
	}
	return n, n, err
}

func (unicodeFoldTransformer) Reset() {}
