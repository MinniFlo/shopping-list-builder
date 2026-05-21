package main

import (
	"charm.land/bubbles/v2/key"
)

type key_map struct {
	Quit        key.Binding
	Down        key.Binding
	Up          key.Binding
	SectionDown key.Binding
	SectionUp   key.Binding
	SetCategory key.Binding
	Increment   key.Binding
	Decrement   key.Binding
	Left        key.Binding
	Right       key.Binding
	Confirm     key.Binding
	Help        key.Binding
}

func (k key_map) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k key_map) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.SectionDown, k.SectionUp},
		{k.Increment, k.Decrement, k.SetCategory, k.Confirm, k.Quit},
	}
}

var default_key_map = key_map{
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move up"),
	),
	SectionDown: key.NewBinding(
		key.WithKeys("tab", "ctrl+d"),
		key.WithHelp("tab/ctrl+d", "next section"),
	),
	SectionUp: key.NewBinding(
		key.WithKeys("shift+tab", "ctrl+u"),
		key.WithHelp("󰘶+tab/ctrl+u", "previous section"),
	),
	SetCategory: key.NewBinding(
		key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9"),
		key.WithHelp("1-9", "1:🌱 2:🥕 3:🧀 4:🍙 5:🧊 6:🍝 7:🥛 8:🍫 9:' '"),
	),
	Increment: key.NewBinding(
		key.WithKeys("+"),
		key.WithHelp("+", "increment section"),
	),
	Decrement: key.NewBinding(
		key.WithKeys("-"),
		key.WithHelp("-", "decrement section"),
	),
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("←/→/h/l", "set staging"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("→/l", "set staging"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "save list to file"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "show key bindings"),
	),
}
