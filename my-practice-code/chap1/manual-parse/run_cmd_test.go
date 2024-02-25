package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		c      config
		input  string
		output string
		err    error
	}{
		{
			c:      config{printUsage: true},
			output: usageString,
		},
		{
			c:      config{numTimes: 5},
			input:  "",
			output: strings.Repeat("Your name please? Press the Enter key when done.\n", 1), // 重复1次
			err:    errors.New("name cannot be empty"),
		},
		{
			c:      config{numTimes: 2},
			input:  "leo\n",
			output: "Your name please? Press the Enter key when done.\n" + strings.Repeat("Hello, leo!\n", 2),
			err:    nil,
		},
	}

	byteBuf := new(bytes.Buffer) // 创建一个新的Buffer, 用于存储runCmd输出的内容
	for _, tc := range tests {
		r := strings.NewReader(tc.input) // 创建一个新的Reader
		err := runCmd(r, byteBuf, tc.c)
		if err != nil && tc.err == nil {
			t.Fatalf("Expected nil error, got: %v\n", err)
		}

		if tc.err != nil {
			if err.Error() != tc.err.Error() {
				t.Fatalf("Expected error: %v, Got error: %v\n", tc.err.Error(), err.Error())
			}
		}

		gotMsg := byteBuf.String() // 获取runCmd的Buffer中的内容
		if gotMsg != tc.output {
			t.Errorf("Expected stdout message to be: %v, Got: %v\n", tc.output, gotMsg)
		}

		byteBuf.Reset() // 重置Buffer
	}
}
