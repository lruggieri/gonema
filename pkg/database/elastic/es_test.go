package elastic

import (
	"testing"
)

func TestGetStringGrams(t *testing.T){
	type test struct{
		initialString string
		separator string
		expectedNGrams map[string]struct{}
		minNGrams int
	}

	tests := []test{
		{
			initialString:"Hello",
			separator:" ",
			expectedNGrams: map[string]struct{}{
				"Hello": struct{}{},
			},
			minNGrams:1,
		},
		{
			initialString:"Hello World",
			separator:" ",
			expectedNGrams: map[string]struct{}{
				"Hello": struct{}{},
				"World": struct{}{},
				"Hello World": struct{}{},
				"World Hello": struct{}{},
			},
			minNGrams:1,
		},
		{
			initialString:"Hello",
			separator:" ",
			expectedNGrams: map[string]struct{}{
				"Hello": struct{}{},
			},
			minNGrams:1,
		},
		{
			initialString:"I am me",
			separator:" ",
			expectedNGrams: map[string]struct{}{
				"I": struct{}{},
				"am": struct{}{},
				"me": struct{}{},
				"I am": struct{}{},
				"am I": struct{}{},
				"I me": struct{}{},
				"me I": struct{}{},
				"am me": struct{}{},
				"me am": struct{}{},
				"I am me": struct{}{},
				"I me am": struct{}{},
				"me I am": struct{}{},
				"me am I": struct{}{},
				"am I me": struct{}{},
				"am me I": struct{}{},
			},
			minNGrams:1,
		},
		{// as before, but changing minNGrams to 2
			initialString:"I am me",
			separator:" ",
			expectedNGrams: map[string]struct{}{
				"I am": struct{}{},
				"am I": struct{}{},
				"I me": struct{}{},
				"me I": struct{}{},
				"am me": struct{}{},
				"me am": struct{}{},
				"I am me": struct{}{},
				"I me am": struct{}{},
				"me I am": struct{}{},
				"me am I": struct{}{},
				"am I me": struct{}{},
				"am me I": struct{}{},
			},
			minNGrams:2,
		},
		{// as before, but changing minNGrams to 3
			initialString:"I am me",
			separator:" ",
			expectedNGrams: map[string]struct{}{
				"I am me": struct{}{},
				"I me am": struct{}{},
				"me I am": struct{}{},
				"me am I": struct{}{},
				"am I me": struct{}{},
				"am me I": struct{}{},
			},
			minNGrams:3,
		},
		{// as before, but changing minNGrams to 4
			initialString:"I am me",
			separator:" ",
			expectedNGrams: map[string]struct{}{},
			minNGrams:4,
		},
	}

	for tstIdx, tst := range tests{
		res := getStringGrams(tst.initialString,tst.separator,tst.minNGrams)
		resMap := make(map[string]struct{},len(res))
		for _,ngram := range res{
			resMap[ngram] = struct{}{}
		}

		if len(tst.expectedNGrams) != len(resMap){
			t.Error("test #",tstIdx,": expected to find ",len(tst.expectedNGrams),"ngram but found ",len(resMap))
			t.FailNow()
		}

		for expectedNgram := range tst.expectedNGrams{
			if _,ok := resMap[expectedNgram] ; !ok{
				t.Error("test #",tstIdx,": expected to find ngram '"+expectedNgram+"' but it was not found")
				t.FailNow()
			}
		}

	}
}
