package match

type Nothing struct{}

func NewNothing() Nothing {
	return Nothing{}
}

func (n Nothing) Match(s string) bool {
	return len(s) == 0
}

func (n Nothing) Index(string) (int, []int) {
	return 0, segments0
}

func (n Nothing) Len() int {
	return lenZero
}

func (n Nothing) String() string {
	return "<nothing>"
}
