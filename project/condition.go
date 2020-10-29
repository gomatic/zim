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
	"context"
	"fmt"
)

// CheckConditions returns true if the Rule should execute based on all of its
// conditions being met. The provided executor is used to run any scripting
// required to check the conditions.
func CheckConditions(
	ctx context.Context,
	r *Rule,
	opts RunOpts,
	executor Executor,
) (bool, error) {

	if !r.when.IsEmpty() {
		// A when condition is defined
		whenCondition, err := CheckCondition(ctx, r, r.when, opts, executor)
		if err != nil {
			return false, err
		}
		if !whenCondition {
			// The "when" condition evaluted to false: condition not met
			return false, nil
		}
	}

	if !r.unless.IsEmpty() {
		// An unless condition is defined
		unlessCondition, err := CheckCondition(ctx, r, r.unless, opts, executor)
		if err != nil {
			return false, err
		}
		if unlessCondition {
			// The "unless" condition evaluted to true: condition not met
			return false, nil
		}
	}

	// All conditions met
	return true, nil
}

// CheckCondition returns true if the given Rule condition is met. The provided
// executor is used to run any scripting required to check the conditions.
func CheckCondition(
	ctx context.Context,
	r *Rule,
	c Condition,
	opts RunOpts,
	executor Executor,
) (bool, error) {

	if c.ResourceExists != "" {
		// The "resource exists" condition evaluates to true if one or more resources
		// match the provided filename or glob pattern
		resources, err := matchResources(r.Component(), r.inProvider, []string{c.ResourceExists})
		if err != nil {
			return false, err
		}
		return len(resources) > 0, nil
	}

	if c.ScriptSucceeds != "" {
		// The "script succeeds" condition evaluates to true if the specified shell
		// command exits without error when run in bash
		err := executor.Execute(ctx, ExecOpts{
			Command:          c.ScriptSucceeds,
			WorkingDirectory: r.Component().Directory(),
			Env:              flattenEnvironment(r.BaseEnvironment()),
			Image:            r.Image(),
			Name:             fmt.Sprintf("%s.condition", r.NodeID()),
			Stdout:           opts.Output,
			Stderr:           opts.Output,
			Debug:            opts.Debug,
			Cmdout:           opts.DebugOutput,
		})
		return err == nil, nil
	}
	return true, nil
}