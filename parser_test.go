package pgot

import (
	"os"
	"testing"
)

const fileName = "./data/daron_mac.pdf"

func TestParserGetStartXref(t *testing.T) {
	f, err := os.Open(fileName)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	p := NewParser()
	sXref, err := p.getStartXref(f)
	if err != nil {
		t.Error(err)
	}
	if sXref != 63544 {
		t.Error("startxref should be 63544")
	}
}

func TestParserparseXref(t *testing.T) {
	f, _ := os.Open(fileName)
	defer f.Close()
	p := NewParser()
	sXreft, _ := p.getStartXref(f)
	err := p.parseXref(f, sXreft)
	if err != nil {
		t.Error(err)
	}

	if p.xrefList.start != 0 {
		t.Error("xref start should be 0")
	}
	if p.xrefList.num != 51 {
		t.Error("xref number should be 51")
	}
	if p.xrefList.trailer.id[0] != `b1ac0eca5f7b1ee4387dc847913d9a93` &&
		p.xrefList.trailer.id[1] != `b1ac0eca5f7b1ee4387dc847913d9a93` {
		t.Error("tralier id should be b1ac0eca5f7b1ee4387dc847913d9a93 and b1ac0eca5f7b1ee4387dc847913d9a93")
	}
	if p.xrefList.trailer.info != `1_0` {
		t.Error("tralier info should be 1_0")
	}
	if p.xrefList.trailer.size != `51` {
		t.Error("tralier size should be 51")
	}
	if p.xrefList.trailer.root != `25_0` {
		t.Error("tralier root should be 25_0")
	}
}
