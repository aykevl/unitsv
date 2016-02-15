// writer.go
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
)

type Writer struct {
	writer       *bufio.Writer
	headerLength int
}

// NewWriter writes the header and returns a new Writer.
func NewWriter(out *bufio.Writer, columns []string) (*Writer, error) {
	w := &Writer{
		writer:       out,
		headerLength: len(columns),
	}

	err := w.WriteRow(columns)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// WriteRow writes a singe row of fields to the file. It checks whether the row
// has the right size.
func (w *Writer) WriteRow(row []string) error {
	if len(row) != w.headerLength {
		return ErrInvalidRowLength
	}

	// It might be possible to optimize this a lot, but that's not yet relevant
	// for me.
	for i, field := range row {
		if i > 0 {
			w.writer.WriteString("\t")
		}

		for _, c := range field {
			switch c {
			case '\n':
				w.writer.WriteString("\\n")
			case '\t':
				w.writer.WriteString("\\t")
			case '\\':
				w.writer.WriteString("\\\\")
			default:
				w.writer.WriteRune(c)
			}
		}
	}

	// This will catch any errors: bufio.Writer ensures that subsequent writes
	// after an error will return that error.
	_, err := w.writer.WriteString("\n")
	return err
}

// Flush flushes the underlying bufio.Writer. This is only necessary if you
// don't flush it yourself.
func (w *Writer) Flush() error {
	return w.writer.Flush()
}
