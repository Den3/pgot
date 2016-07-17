package pgot

type Xref struct {
	Offset uint64

	//Generation Number
	GenNum uint16

	// n for an in-use entry or f for a free entry
	Entry byte
}

type Trailer struct {
	Size string
	Root string
	Info string
	ID   [2]string
}

type XrefList struct {
	Start   uint16
	Num     uint32
	List    []Xref
	Trailer Trailer
}
