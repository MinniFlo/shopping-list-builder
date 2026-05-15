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
