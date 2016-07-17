package pgot

type xref struct {
	offset uint64

	//Generation Number
	genNum uint16

	// n for an in-use entry or f for a free entry
	entry byte
}

type trailer struct {
	size string
	root string
	info string
	id   [2]string
}

type xrefList struct {
	start   uint16
	num     uint32
	list    []xref
	trailer trailer
}
