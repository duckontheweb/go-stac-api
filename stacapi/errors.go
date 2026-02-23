package stacapi

import "fmt"

type CollectionNotFoundError struct {
	Id string
}

func (e CollectionNotFoundError) Error() string {
	return fmt.Sprintf("Collection with ID '%s' not found.", e.Id)
}
