package api

import (
	"fmt"
	"time"
)

type Visitor struct {
	Recursive bool
	Rate      time.Duration
}

func NewVisitor(recursive bool, rate int) (*Visitor, error) {
	v := &Visitor{
		Recursive: recursive,
		Rate:      time.Duration(rate) * time.Millisecond,
	}
	return v, nil
}

func (v *Visitor) VisitTree() {
	fmt.Println("Visiting GitHub...")
}
