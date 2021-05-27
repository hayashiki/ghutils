package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/jszwec/csvutil"
)

type csvStruct struct{}

//go:generate mockgen -source ./csv/csv.go -package mock_csv -destination ./csv/mock/csv.go
type CSV interface {
	Generate(items interface{}) (io.ReadSeeker, error)
	GenerateBytes(items interface{}, withHeader bool) (bytes.Buffer, error)
}

func New() CSV {
	return &csvStruct{}
}

func (c *csvStruct) Generate(items interface{}) (io.ReadSeeker, error) {
	b, err := csvutil.Marshal(items)
	if err != nil {
		fmt.Printf("CSV Marshal error=%v", err)
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func (c *csvStruct) GenerateBytes(items interface{}, withHeader bool) (bytes.Buffer, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.Comma = '\t'
	enc := csvutil.NewEncoder(w)
	enc.AutoHeader = withHeader
	if err := enc.Encode(items); err != nil {
		fmt.Println("error:", err)
	}

	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Println("error:", err)
	}

	return buf, nil
}
