package ydb

import (
	"context"
	"fmt"
	"strings"
)

type (
	Logger struct {
		indent int
	}
	LoggerEntry struct {
		logger       *Logger
		functionName string
		indent       int
	}
)

func (l *Logger) Start(functionName string, args ...any) *LoggerEntry {
	if l == nil {
		return nil
	}

	l.indent++

	fmt.Printf("[YDB]%s-> %s(", strings.Repeat("  ", l.indent), functionName)
	for i, arg := range args {
		if i > 0 {
			fmt.Printf(", ")
		}
		if _, has := arg.(context.Context); has {
			arg = "ctx"
		}
		fmt.Printf("%v", arg)
	}
	fmt.Printf(")\n")

	return &LoggerEntry{
		logger:       l,
		functionName: functionName,
		indent:       l.indent,
	}
}

func (l *LoggerEntry) End(results ...any) {
	if l == nil {
		return
	}

	defer func() {
		l.logger.indent--
	}()

	fmt.Printf("[YDB]%s<- %s => ", strings.Repeat("  ", l.indent), l.functionName)
	for i, res := range results {
		if i > 0 {
			fmt.Printf(", ")
		}
		if v, has := res.([]byte); has {
			res = v
		}
		fmt.Printf("%v", res)
	}
	fmt.Printf("\n")
}
