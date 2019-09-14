package rarbg

import (
	"github.com/otiai10/gosseract"
	"testing"
)

func TestTesseract_ImageDecoding(t *testing.T){

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage("testing_files/captcha.png")
	text, _ := client.Text()

	if text != "MUWHK"{
		t.Error("tesseract result not matching expected result")
		t.FailNow()
	}

}