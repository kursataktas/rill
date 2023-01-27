package blob

import (
	"context"
	"io"
	"os"
	"strings"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"gocloud.dev/blob"
)

// todo :: check if string conversions can be avoided
func downloadCSV(ctx context.Context, bucket *blob.Bucket, obj *blob.ListObject, option *extractConfig, fw *os.File) error {
	reader := NewBlobObjectReader(ctx, bucket, obj)

	rows, err := csvRows(reader, option)
	if err != nil {
		return err
	}

	// write rows
	for _, r := range rows {
		if _, err := fw.WriteString(r); err != nil {
			return err
		}
		if _, err := fw.WriteString("\n"); err != nil {
			return err
		}
	}

	return nil
}

func csvRows(reader *ObjectReader, option *extractConfig) ([]string, error) {
	if option.strategy == runtimev1.Source_ExtractPolicy_TAIL {
		return csvRowsTail(reader, option)
	}
	return csvRowsHead(reader, option)
}

func csvRowsTail(reader *ObjectReader, option *extractConfig) ([]string, error) {
	header, err := getHeader(reader)
	if err != nil {
		return nil, err
	}

	remBytes := int64(option.limtInBytes - uint64(len([]byte(header))))
	if _, err := reader.Seek(0-remBytes, io.SeekEnd); err != nil {
		return nil, err
	}

	p := make([]byte, remBytes)
	if _, err := reader.Read(p); err != nil {
		return nil, err
	}

	rows := strings.Split(string(p), "\n")
	// remove first row (possibly incomplete)
	rows = rows[1:]
	// append header at start
	return append([]string{header}, rows...), nil
}

func csvRowsHead(reader *ObjectReader, option *extractConfig) ([]string, error) {
	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	p := make([]byte, option.limtInBytes)
	if _, err := reader.Read(p); err != nil {
		return nil, err
	}

	rows := strings.Split(string(p), "\n")
	// remove last row (possibly incomplete)
	return rows[:len(rows)-1], nil
}

// tries to get csv header from reader by incrmentally reading 1KB bytes
func getHeader(r *ObjectReader) (string, error) {
	fetchLength := 1024
	p := make([]byte, 0)
	for {
		temp := make([]byte, fetchLength)
		n, err := r.Read(temp)
		if err != nil && !strings.Contains(err.Error(), "EOF") {
			return "", err
		}

		p = append(p, temp...)
		rows := strings.Split(string(p), "\n")
		if len(rows) > 1 {
			// complete header found
			return rows[0], nil
		}

		if n < fetchLength {
			// end of csv
			return "", io.EOF
		}
	}
}
