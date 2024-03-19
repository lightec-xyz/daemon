package circuits

type Recursive struct {
}

func NewRecursive() *Recursive {
	return &Recursive{}
}

func (r *Recursive) Verify(opt *OptRecursive) (bool, error) {
	panic(opt)
}

func (r *Recursive) GenerateProof(opt *OptRecursive) error {
	panic(opt)
}

type OptRecursive struct {
}
