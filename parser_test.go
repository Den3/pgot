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

	if p.XrefList.Start != 0 {
		t.Error("xref start should be 0")
	}
	if p.XrefList.Num != 51 {
		t.Error("xref number should be 51")
	}
	if p.XrefList.Trailer.ID[0] != `b1ac0eca5f7b1ee4387dc847913d9a93` &&
		p.XrefList.Trailer.ID[1] != `b1ac0eca5f7b1ee4387dc847913d9a93` {
		t.Error("tralier id should be b1ac0eca5f7b1ee4387dc847913d9a93 and b1ac0eca5f7b1ee4387dc847913d9a93")
	}
	if p.XrefList.Trailer.Info != `1_0` {
		t.Error("tralier info should be 1_0")
	}
	if p.XrefList.Trailer.Size != `51` {
		t.Error("tralier size should be 51")
	}
	if p.XrefList.Trailer.Root != `25_0` {
		t.Error("tralier root should be 25_0")
	}
}
