// Copyright 2023 dvonthenen ChatGPT Proxy contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache License 2.0

package interfaces

import "errors"

type SkillType int64

const (
	SkillTypeDefault   SkillType = iota
	SkillTypeGeneric             = 1
	SkillTypeExpert              = 2
	SkillTypeDAN                 = 991
	SkillTypeSTAN                = 992
	SkillTypeDUDE                = 993
	SkillTypeJailBreak           = 994
	SkillTypeMongo               = 995
)

var (
	// ErrInvalidInput required input was not found
	ErrInvalidInput = errors.New("required input was not found")
)
