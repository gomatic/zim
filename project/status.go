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

// Code indicates the scheduling result for a Rule
type Code int

const (

	// Error indicates the Rule could not be run
	Error Code = iota

	// UpToDate indicates the Rule is up-to-date and doesn't need to be run
	UpToDate

	// ExecError indicates Rule execution was attempted but failed
	ExecError

	// MissingOutputError indicates a Rule output was not produced
	MissingOutputError

	// OK indicates that the Rule executed successfully
	OK

	// Cached indicates the Rule artifact was cached
	Cached
)
