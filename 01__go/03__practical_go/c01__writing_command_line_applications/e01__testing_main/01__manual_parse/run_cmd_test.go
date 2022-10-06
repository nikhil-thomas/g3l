package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		c             config
		input, output string
		err           error
	}{
		{
			c:      config{printUsage: true},
			output: usage,
		},
		{
			c:      config{numTimes: 3},
			input:  "",
			output: strings.Repeat("Your name please? Press the Enter key when done.\n", 1),
			err:    errors.New("you didn't enter your name"),
		},
		{
			c:      config{numTimes: 5},
			input:  "N T",
			output: "Your name please? Press the Enter key when done.\n" + strings.Repeat("Nice to meet you N T\n", 5),
			err:    nil,
		},
	}
	byteBuf := new(bytes.Buffer)
	for _, tc := range tests {
		r := strings.NewReader(tc.input)
		err := runCmd(r, byteBuf, tc.c)
		if err != nil && tc.err == nil {
			t.Fatalf("Expected nil error, got: %v\n", err)
		}
		if tc.err != nil {
			if err.Error() != tc.err.Error() {
				t.Fatalf("Expected error: %v, Got error: %v\n", tc.err.Error(), err.Error())
			}
		}
		gotMsg := byteBuf.String()
		if gotMsg != tc.output {
			t.Errorf("Expected stdout mesage to be: %v, Got: %v\n", tc.output, gotMsg)
		}
		byteBuf.Reset()

	}
}
