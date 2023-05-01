package plug

import "fmt"

type Plug interface {
	Off()
	On()
	IDs()
}

type MyPlug struct {
	state string
}

func (p *MyPlug) Off(id string) {
	fmt.Println(p.state, id)
}
