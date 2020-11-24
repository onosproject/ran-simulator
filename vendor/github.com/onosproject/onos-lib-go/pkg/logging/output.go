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
	"bytes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/url"
	"strings"
	"sync"
)

func newZapOutput(logger LoggerConfig, output OutputConfig, sink SinkConfig) (*zapOutput, error) {
	zapConfig := zap.Config{}
	zapConfig.Level = levelToAtomicLevel(output.GetLevel())
	zapConfig.Encoding = string(sink.GetEncoding())
	zapConfig.EncoderConfig.EncodeName = zapcore.FullNameEncoder
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.EncoderConfig.EncodeDuration = zapcore.NanosDurationEncoder
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	zapConfig.EncoderConfig.NameKey = "logger"
	zapConfig.EncoderConfig.MessageKey = "message"
	zapConfig.EncoderConfig.LevelKey = "level"
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.CallerKey = "caller"
	zapConfig.EncoderConfig.StacktraceKey = "trace"

	var encoder zapcore.Encoder
	switch sink.GetEncoding() {
	case ConsoleEncoding:
		encoder = zapcore.NewConsoleEncoder(zapConfig.EncoderConfig)
	case JSONEncoding:
		encoder = zapcore.NewJSONEncoder(zapConfig.EncoderConfig)
	}

	var path string
	switch sink.GetType() {
	case StdoutSinkType:
		path = StdoutSinkType.String()
	case StderrSinkType:
		path = StderrSinkType.String()
	case FileSinkType:
		path = sink.GetFileSinkConfig().Path
	case KafkaSinkType:
		kafkaConfig := sink.GetKafkaSinkConfig()
		var rawQuery bytes.Buffer
		if kafkaConfig.Topic != "" {
			rawQuery.WriteString("topic=")
			rawQuery.WriteString(kafkaConfig.Topic)
		}

		if kafkaConfig.Key != "" {
			rawQuery.WriteString("&")
			rawQuery.WriteString("key=")
			rawQuery.WriteString(kafkaConfig.Key)
		}
		kafkaURL := url.URL{Scheme: KafkaSinkType.String(), Host: strings.Join(kafkaConfig.Brokers, ","), RawQuery: rawQuery.String()}
		path = kafkaURL.String()
	}

	writer, err := getWriter(path)
	if err != nil {
		return nil, err
	}

	atomLevel := zap.AtomicLevel{}
	switch output.GetLevel() {
	case DebugLevel:
		atomLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case InfoLevel:
		atomLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case WarnLevel:
		atomLevel = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case ErrorLevel:
		atomLevel = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case PanicLevel:
		atomLevel = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	case FatalLevel:
		atomLevel = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	}

	zapLogger, err := zapConfig.Build(zap.AddCallerSkip(2))
	if err != nil {
		return nil, err
	}

	zapLogger = zapLogger.WithOptions(
		zap.WrapCore(
			func(zapcore.Core) zapcore.Core {
				return zapcore.NewCore(encoder, writer, &atomLevel)
			}))
	return &zapOutput{
		config: output,
		logger: zapLogger.Named(logger.Name),
	}, nil
}

var writers = make(map[string]zapcore.WriteSyncer)
var writersMu = &sync.Mutex{}

func getWriter(url string) (zapcore.WriteSyncer, error) {
	writersMu.Lock()
	defer writersMu.Unlock()
	writer, ok := writers[url]
	if !ok {
		ws, _, err := zap.Open(url)
		if err != nil {
			return nil, err
		}
		writer = ws
		writers[url] = writer
	}
	return writer, nil
}

// Output is a logging output
type Output interface {
	Debug(...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Info(...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})

	Error(...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})

	Fatal(...interface{})
	Fatalf(template string, args ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	Panic(...interface{})
	Panicf(template string, args ...interface{})
	Panicw(msg string, keysAndValues ...interface{})

	DPanic(...interface{})
	DPanicf(template string, args ...interface{})
	DPanicw(msg string, keysAndValues ...interface{})

	Warn(...interface{})
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
}

// zapOutput is a logging output implementation
type zapOutput struct {
	config OutputConfig
	logger *zap.Logger
}

func (o *zapOutput) Debug(args ...interface{}) {
	o.logger.Sugar().Debug(args...)
}

func (o *zapOutput) Debugf(template string, args ...interface{}) {
	o.logger.Sugar().Debugf(template, args...)
}

func (o *zapOutput) Debugw(msg string, keysAndValues ...interface{}) {
	o.logger.Sugar().Debugw(msg, keysAndValues...)
}

func (o *zapOutput) Info(args ...interface{}) {
	o.logger.Sugar().Info(args...)
}

func (o *zapOutput) Infof(template string, args ...interface{}) {
	o.logger.Sugar().Infof(template, args...)
}

func (o *zapOutput) Infow(msg string, keysAndValues ...interface{}) {
	o.logger.Sugar().Infow(msg, keysAndValues...)
}

func (o *zapOutput) Error(args ...interface{}) {
	o.logger.Sugar().Error(args...)
}

func (o *zapOutput) Errorf(template string, args ...interface{}) {
	o.logger.Sugar().Errorf(template, args...)
}

func (o *zapOutput) Errorw(msg string, keysAndValues ...interface{}) {
	o.logger.Sugar().Errorw(msg, keysAndValues...)
}

func (o *zapOutput) Fatal(args ...interface{}) {
	o.logger.Sugar().Fatal(args...)
}

func (o *zapOutput) Fatalf(template string, args ...interface{}) {
	o.logger.Sugar().Fatalf(template, args)
}

func (o *zapOutput) Fatalw(msg string, keysAndValues ...interface{}) {
	o.logger.Sugar().Fatalw(msg, keysAndValues...)
}

func (o *zapOutput) Panic(args ...interface{}) {
	o.logger.Sugar().Panic(args...)
}

func (o *zapOutput) Panicf(template string, args ...interface{}) {
	o.logger.Sugar().Panicf(template, args...)
}

func (o *zapOutput) Panicw(msg string, keysAndValues ...interface{}) {
	o.logger.Sugar().Panicw(msg, keysAndValues...)
}

func (o *zapOutput) DPanic(args ...interface{}) {
	o.logger.Sugar().DPanic(args...)
}

func (o *zapOutput) DPanicf(template string, args ...interface{}) {
	o.logger.Sugar().DPanicf(template, args...)
}

func (o *zapOutput) DPanicw(msg string, keysAndValues ...interface{}) {
	o.logger.Sugar().DPanicw(msg, keysAndValues...)
}

func (o *zapOutput) Warn(args ...interface{}) {
	o.logger.Sugar().Warn(args...)
}

func (o *zapOutput) Warnf(template string, args ...interface{}) {
	o.logger.Sugar().Warnf(template, args...)
}

func (o *zapOutput) Warnw(msg string, keysAndValues ...interface{}) {
	o.logger.Sugar().Warnw(msg, keysAndValues...)
}

var _ Output = &zapOutput{}
