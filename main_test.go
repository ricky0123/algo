package main

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAliases(t *testing.T) {
	tests := []struct {
		Name        string
		Config      string
		TestCommand string
		Expected    string
		BashSetup   string
	}{
		{
			Name:        "one line",
			Config:      "e: echo",
			TestCommand: "e hello world",
			Expected:    "hello world\n",
		},
		{
			Name:        "with newlines",
			Config:      "\ne: echo\n",
			TestCommand: "e hello world",
			Expected:    "hello world\n",
		},
		{
			Name:        "empty",
			Config:      "",
			TestCommand: "",
			Expected:    "",
		},
		{
			Name: "$ syntax",
			Config: `
e:
    $: echo hello
`,
			TestCommand: "e world",
			Expected:    "hello world\n",
		},
		{
			Name:      "multiple word parsing",
			BashSetup: "countArgs() { echo $#; }",
			Config: `
count: countArgs 1 "2 3"
`,
			TestCommand: "count",
			Expected:    "2\n",
		},
		{
			Name: "multiple aliases",
			Config: `
x: echo 1
y: echo 2
`,
			TestCommand: "x; y",
			Expected:    "1\n2\n",
		},
		{
			Name: "progressively building command",
			Config: `
a:
  $: echo
  b:
    $: hello
    c: world
`,
			TestCommand: "a b c",
			Expected:    "hello world\n",
		},
		{
			Name: "multiple subcommands",
			Config: `
a:
  b1: echo you chose b1
  b2: echo you chose b2
  b3: echo you chose b3
`,
			TestCommand: "a b1; a b3; a",
			Expected:    "you chose b1\nyou chose b3\n",
		},
		{
			Name: "with &&",
			Config: `
success: true && echo success
`,
			TestCommand: "success",
			Expected:    "success\n",
		},
		{
			Name:      "with environment variable",
			BashSetup: "export MYVAR=hello",
			Config: `
what: echo $MYVAR
`,
			TestCommand: "what world",
			Expected:    "hello world\n",
		},
		{
			Name:      "with $()",
			BashSetup: "export MYVAR=hello; countArgs() { echo $#; }",
			Config: `
e: echo $(echo 1 2 3)
`,
			TestCommand: "e",
			Expected:    "1 2 3\n",
		},
		{
			Name:      "count extra args correctly",
			BashSetup: "countArgs() { echo $#; }",
			Config: `
what: countArgs
`,
			TestCommand: "what '1 2' \"3 4\" 5",
			Expected:    "3\n",
		},
		{
			Name: "comments",
			Config: `
# top comment
x: echo hello # this is a comment
# again
y: echo world # here is another comment
`,
			TestCommand: "x; y",
			Expected:    "hello\nworld\n",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			source := GenerateSource([]byte(test.Config))
			cmd := exec.Command("bash", "-c", test.BashSetup+"\n"+source+"\n"+test.TestCommand)
			outBytes, err := cmd.CombinedOutput()
			out := string(outBytes)
			require.NoError(t, err, out)
			assert.Equal(t, test.Expected, out)
		})
	}
}
