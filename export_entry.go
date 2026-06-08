package main

type export_entry struct {
	name           string
	amount         float64
	unit           unit
	state          stageing_state
	list_entry_ids []int
	selected       bool
}

type ByListEntryIds []*export_entry

func (e ByListEntryIds) Len() int {
	return len(e)
}

func (e ByListEntryIds) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e ByListEntryIds) Less(i, j int) bool {
	return e[i].list_entry_ids[0] < e[j].list_entry_ids[0]
}
