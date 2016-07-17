package pgot

import (
	"bytes"
	"os"
	"regexp"
	"strconv"
)

type Parser struct {
	XrefList       XrefList
	xrefReg        *regexp.Regexp
	tralierSizeReg *regexp.Regexp
	tralierRootReg *regexp.Regexp
	tralierInfoReg *regexp.Regexp
	tralierIDReg   *regexp.Regexp
}

func NewParser() *Parser {
	return &Parser{
		xrefReg:        regexp.MustCompile(`(\d{10})\s(\d{5})\s([f|n])`),
		tralierSizeReg: regexp.MustCompile(`Size\s+(\d+)`),
		tralierRootReg: regexp.MustCompile(`Root\s+(\d+)\s+(\d+)\s+R`),
		tralierInfoReg: regexp.MustCompile(`Info\s+(\d+)\s+(\d+)\s+R`),
		tralierIDReg:   regexp.MustCompile(`ID\s+\[\s+<(\w+)>\s*<(\w+)>`),
	}
}

func (p *Parser) checkWhiteSpace(i byte) bool {
	// 0x00 null (NUL), 0x09 horizontal tab (HT), 0x0A line feed (LF),
	// 0x0C form feed (FF), 0x0D carriage return (CR), 0x20 space (SP)
	if i == 0x00 || i == 0x09 ||
		i == 0x0A || i == 0x0C ||
		i == 0x0D {
		return true
	}
	return false
}

func (p *Parser) GetStartXref(file *os.File) (int, error) {
	state, err := file.Stat()
	if err != nil {
		return 0, err
	}

	// just guess where is startxref, as spec's description is 1024
	buf := make([]byte, 1024)
	_, err = file.ReadAt(buf, state.Size()-1024)
	if err != nil {
		return 0, err
	}

	sXref := ""
	sXrefIndex := bytes.Index(buf, []byte("startxref"))
	xrefOffset := 0
	if sXrefIndex != -1 {
		sXrefIndex += 9
		for {
			if buf[sXrefIndex] == '%' {
				xrefOffset, err = strconv.Atoi(sXref)
				if err != nil {
					return 0, err
				}
				break
			}
			if buf[sXrefIndex] < 48 || buf[sXrefIndex] > 57 {
				sXrefIndex++
				continue
			}

			sXref += string(buf[sXrefIndex])
			sXrefIndex++
		}
	}
	return xrefOffset, nil
}

func (p *Parser) GetXref(file *os.File, offset int) error {
	buf := make([]byte, 2*1024) // 2M
	file.ReadAt(buf, int64(offset))
	tmpOffset := 4 // len("xref")
	tmpString := ""
	for {
		if buf[tmpOffset] == ' ' {
			tmpOffset++
			break
		}
		if buf[tmpOffset] != '\r' && buf[tmpOffset] != '\n' {
			tmpString += string(buf[tmpOffset])
		}
		tmpOffset++
	}
	val, err := strconv.Atoi(tmpString)
	if err != nil {
		return err
	}
	p.XrefList.Start = uint16(val)

	tmpString = ""
	for {
		if p.checkWhiteSpace(buf[tmpOffset]) {
			break
		}

		tmpString += string(buf[tmpOffset])
		tmpOffset++
	}
	val, err = strconv.Atoi(tmpString)
	if err != nil {
		return err
	}
	p.XrefList.Num = uint32(val)

	xrefList := make([]Xref, val)
	matches := p.xrefReg.FindAllStringSubmatch(string(buf), -1)
	for i := 0; i < val; i++ {
		matchOffset, err := strconv.Atoi(matches[i][1])
		if err != nil {
			return err
		}
		xrefList[i].Offset = uint64(matchOffset)
		matchGenNum, err := strconv.Atoi(matches[i][2])
		if err != nil {
			return err
		}
		xrefList[i].GenNum = uint16(matchGenNum)
		xrefList[i].Entry = byte(matches[i][3][0])
	}
	p.XrefList.List = xrefList

	// tralier
	size := p.tralierSizeReg.FindStringSubmatch(string(buf))
	p.XrefList.Trailer.Size = size[1]
	root := p.tralierRootReg.FindStringSubmatch(string(buf))
	p.XrefList.Trailer.Root = root[1] + "_" + root[2]
	info := p.tralierInfoReg.FindStringSubmatch(string(buf))
	p.XrefList.Trailer.Info = info[1] + "_" + info[2]
	id := p.tralierIDReg.FindStringSubmatch(string(buf))
	p.XrefList.Trailer.ID[0] = id[1]
	p.XrefList.Trailer.ID[1] = id[2]
	return nil
}
