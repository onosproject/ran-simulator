// Copyright 2020-present Open Networking Foundation.
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

package logging

import (
	zp "go.uber.org/zap"
	zc "go.uber.org/zap/zapcore"
	"strings"
)

// Level :
type Level int

const (
	// DebugLevel logs a message at debug level
	DebugLevel Level = iota
	// InfoLevel logs a message at info level
	InfoLevel
	// WarnLevel logs a message at warning level
	WarnLevel
	// ErrorLevel logs a message at error level
	ErrorLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// DPanicLevel logs at PanicLevel; otherwise, it logs at ErrorLevel
	DPanicLevel

	// EmptyLevel :
	EmptyLevel
)

// String :
func (l Level) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", "PANIC", "DPANIC", ""}[l]
}

func levelToAtomicLevel(l Level) zp.AtomicLevel {
	switch l {
	case DebugLevel:
		return zp.NewAtomicLevelAt(zc.DebugLevel)
	case InfoLevel:
		return zp.NewAtomicLevelAt(zc.InfoLevel)
	case WarnLevel:
		return zp.NewAtomicLevelAt(zc.WarnLevel)
	case ErrorLevel:
		return zp.NewAtomicLevelAt(zc.ErrorLevel)
	case FatalLevel:
		return zp.NewAtomicLevelAt(zc.FatalLevel)
	case PanicLevel:
		return zp.NewAtomicLevelAt(zc.PanicLevel)
	case DPanicLevel:
		return zp.NewAtomicLevelAt(zc.DPanicLevel)
	}
	return zp.NewAtomicLevelAt(zc.ErrorLevel)
}

func levelStringToLevel(l string) Level {
	switch strings.ToUpper(l) {
	case DebugLevel.String():
		return DebugLevel
	case InfoLevel.String():
		return InfoLevel
	case WarnLevel.String():
		return WarnLevel
	case ErrorLevel.String():
		return ErrorLevel
	case FatalLevel.String():
		return FatalLevel
	case PanicLevel.String():
		return PanicLevel
	case DPanicLevel.String():
		return DPanicLevel
	}
	return ErrorLevel
}
