package model

import (
	"fmt"
)

type sql struct {
	rootPath string
}

func (s *sql) Generate(m *model) error {
	fmt.Println("CREATE SQL Script")
	return nil
}
