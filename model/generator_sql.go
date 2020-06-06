package model

import (
	"fmt"
)

type SQL struct{}

func (s *SQL) Generate(m Model) error {
	fmt.Println("CREATE SQL Script")
	return nil
}
