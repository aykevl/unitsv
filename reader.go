// reader.go
//
// Copyright (c) 2016, Ayke van Laethem
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS
// IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED
// TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A
// PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package unitsv

import (
	"bufio"
	"bytes"
	"io"
	"net/textproto"
	"strings"
)

// Reader implements the TSV reader.
type Reader struct {
	reader       *textproto.Reader
	columnMap    map[int]int // ReadRow() index mapped to file row index
	headerLength int
	rowLength    int
}

// Config should be passed to NewReader with at least one of Required or
// Optional columns set.
type Config struct {
	Required []string // required columns
	Optional []string // optional columns
}

// NewReader creates a new TSV reader. Calling this function reads the header
// and verifies it with the 'columns' parameter.
func NewReader(in io.Reader, config Config) (*Reader, error) {
	reader := textproto.NewReader(bufio.NewReader(in))

	header, err := reader.ReadLine()
	if err != nil && err != io.EOF {
		return nil, err
	}
	columnIndex := make(map[string]int)
	fields, err := splitTsvFields(header)
	if err != nil {
		return nil, err
	}
	for i, column := range fields {
		columnIndex[column] = i
	}

	columnMap := make(map[int]int, len(config.Required)+len(config.Optional))
	for i, column := range config.Required {
		index, ok := columnIndex[column]
		if !ok {
			return nil, ErrColumns
		}
		columnMap[i] = index
	}
	for i, column := range config.Optional {
		if index, ok := columnIndex[column]; ok {
			columnMap[i+len(config.Required)] = index
		}
	}

	r := &Reader{
		reader:       reader,
		columnMap:    columnMap,
		headerLength: len(fields),
		rowLength:    len(config.Required) + len(config.Optional),
	}

	return r, nil
}

// ReadRow reads a single row from the TSV file. The returned slice is the
// length of the required and optional columns combined and has fields at the
// same position as in the column list. The actual columns in the file may be at
// a different index.
func (r *Reader) ReadRow() ([]string, error) {
	line, err := r.reader.ReadLine()
	if err != nil {
		// could also be io.EOF
		return nil, err
	}
	fields, err := splitTsvFields(line)
	if err != nil {
		return nil, err
	}
	if len(fields) != r.headerLength {
		return nil, ErrParsingTSV
	}
	row := make([]string, r.rowLength)
	for i := 0; i < len(r.columnMap); i++ {
		if index, ok := r.columnMap[i]; ok {
			row[i] = fields[index]
		}
	}
	return row, nil
}

// splitTsvFields separates tab-seapareted-values and unescapes them.
func splitTsvFields(line string) ([]string, error) {
	fields := strings.Split(line, "\t")
	for i := 0; i < len(fields); i++ {
		var field bytes.Buffer
		escape := false
		for _, c := range fields[i] {
			if !escape {
				if c == '\\' {
					escape = true
				} else {
					field.WriteRune(rune(c))
				}
			} else {
				switch c {
				case 't':
					field.WriteRune('\t')
				case 'n':
					field.WriteRune('\n')
				case '\\':
					field.WriteRune('\\')
				default:
					return nil, ErrParsingTSV
				}
				escape = false
			}
		}
		if escape {
			return nil, ErrParsingTSV
		}

		fields[i] = field.String()
	}
	return fields, nil
}
