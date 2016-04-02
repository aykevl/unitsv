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
	"bytes"
	"io"
	"testing"
)

var testInput = `header	header2	unused	abc
1	a	-	3002234232222342
u		-	x
ü	...	-	n
\t	\n	-	\\n
`
var testHeadersReadRequired = []string{"header2", "header"}
var testHeadersReadOptional = []string{"abc", "noheader"}
var testHeadersWrite = []string{"header2", "header", "abc"}
var testData = [][]string{
	{"a", "1", "3002234232222342"},
	{"", "u", "x"},
	{"...", "ü", "n"},
	{"\n", "\t", "\\n"},
}
var testOutput = `header2	header	abc
a	1	3002234232222342
	u	x
...	ü	n
\n	\t	\\n
`

func TestReader(t *testing.T) {
	reader, err := NewReader(bytes.NewBufferString(testInput), Config{
		Required: testHeadersReadRequired,
		Optional: testHeadersReadOptional,
	})
	if err != nil {
		t.Fatal("error while opening reader:", err)
	}
	for i_row, row := range testData {
		rowParsed, err := reader.ReadRow()
		if err != nil {
			t.Fatal("error while reading row:", err)
		}
		if len(rowParsed) != len(row)+1 {
			t.Errorf("row length is not equal to header length + 1 for row %d", i_row)
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

func TestWriter(t *testing.T) {
	outfile := &bytes.Buffer{}
	writer, err := NewWriter(outfile, testHeadersWrite)
	if err != nil {
		t.Fatal("error while opening writer:", err)
	}
	for _, row := range testData {
		err := writer.WriteRow(row)
		if err != nil {
			t.Errorf("error writing row %#v: %s", row, err)
		}
	}
	err = writer.Flush()
	if err != nil {
		t.Error("error while flushing:", err)
	}

	output := outfile.Bytes()
	if !bytes.Equal([]byte(testOutput), output) {
		t.Errorf("the expected and actual output does not match, expected and actual:\n%#v\n%#v\n", testOutput, string(output))
	}
}
