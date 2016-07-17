package pgot

import (
	"os"
	"testing"
)

const fileName = "./data/daron_mac.pdf"

// const fileName = "./data/Hire_Chia.pdf"

func TestParserGetStartXref(t *testing.T) {
	f, err := os.Open(fileName)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	p := NewParser()
	sXref, err := p.GetStartXref(f)
	if err != nil {
		t.Error(err)
	}
	if sXref != 63544 {
		t.Error("startxref should be 63544")
	}
}

func TestParserGetXref(t *testing.T) {
	f, _ := os.Open(fileName)
	defer f.Close()
	p := NewParser()
	sXreft, _ := p.GetStartXref(f)
	err := p.GetXref(f, sXreft)
	if err != nil {
		t.Error(err)
	}
	if p.XrefList.Num != 51 {
		t.Error("xref start should be 51")
	}
}
