// unitsv_test.go
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
	"io"
	"strings"
	"testing"
)

var testInput = `header	header2	abc
1	a	3002234232222342
u		x
ü	...	n
\t	\n	\\n
`
var testHeaders = []string{"header2", "header", "abc"}
var testParsed = [][]string{
	{"a", "1", "3002234232222342"},
	{"", "u", "x"},
	{"...", "ü", "n"},
	{"\n", "\t", "\\n"},
}

func TestReader(t *testing.T) {
	infile := strings.NewReader(testInput)
	reader, err := NewReader(bufio.NewReader(infile), testHeaders)
	if err != nil {
		t.Fatal("error while opening reader:", err)
	}
	for i_row, row := range testParsed {
		rowParsed, err := reader.ReadRow()
		if err != nil {
			t.Fatal("error while reading row:", err)
		}
		if len(rowParsed) != len(row) {
			t.Errorf("row length is not equal to header length for row %d", i_row)
			continue
		}
		for i_field, field := range row {
			if field != rowParsed[i_field] {
				t.Errorf("expected and actual value differ (row %d, column %d, expected %#v, actual %#v)", i_row, i_field, field, rowParsed[i_field])
			}
		}
	}
	if _, err := reader.ReadRow(); err != io.EOF {
		t.Error("expected EOF when last row was read")
	}
}
