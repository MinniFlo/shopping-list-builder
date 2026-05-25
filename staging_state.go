package main

type stageing_state int

const (
	NOT_STAGED stageing_state = iota
	MABY
	STAGED
)

const stage_state_max = 2

func (i *stageing_state) Next() {
	if *i < stage_state_max {
		*i += 1
	}
}

func (i *stageing_state) Prev() {
	if *i > 0 {
		*i -= 1
	}
}

func MergeStagingState(s1, s2 stageing_state) stageing_state {
	if s1 == s2 {
		return s1
	} else {
		return MABY
	}
}

func (s stageing_state) String() string {
	return [...]string{" ", "", " "}[s]
}
