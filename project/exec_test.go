// Copyright 2020 Fugue, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package project

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestBashExecutor(t *testing.T) {

	dir := testDir()
	ctx := context.Background()
	e := NewBashExecutor()

	var stdout bytes.Buffer

	err := e.Execute(ctx, ExecOpts{
		Command:          "echo HI $PWD",
		WorkingDirectory: dir,
		Stdout:           &stdout,
	})
	require.Nil(t, err)

	expected := fmt.Sprintf("HI %s", dir)

	// Somehow the tmpdir is prefixed with /private on macos when using PWD
	expectedAlt := fmt.Sprintf("HI /private%s", dir)

	out := strings.TrimSpace(stdout.String())

	if out != expected && out != expectedAlt {
		t.Error("Unexpected output:", out)
	}
}
