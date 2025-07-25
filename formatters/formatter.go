package formatters

import "github.com/eskandaridanial/blink/models"

type Formatter interface {
	Format(record models.Record) []byte
}
