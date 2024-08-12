package dfinity

type BlockHeight struct {
	Hash      string
	Height    uint32
	Signature string
}

type Signature struct {
	Signed    bool
	Signature []string
}
