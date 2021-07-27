package hima

type State struct {
	score int
}

func CreateState() *State {
	return &State{}
}
