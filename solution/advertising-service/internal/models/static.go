package models

import "io"

type Static struct {
	Data        io.Reader
	Size        int64
	ContentType string
}
