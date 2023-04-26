// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package personas

import (
	"errors"
)

const (
	DefaultPort int = 443
)

var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")

	// ErrNotFoundOpenAiKey openai key not found
	ErrNotFoundOpenAiKey = errors.New("openai key not found")

	// ErrInvalidOpenAiClient invalid open ai client
	ErrInvalidOpenAiClient = errors.New("invalid open ai client")
)
