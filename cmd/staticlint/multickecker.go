// Пакет main - главный пакет мультичекера.
//
// Мультичекер состоит из следующих анализаторов:
//  1. Стандартные анализаторы - printf, shadow, assign, defers, nilness, unreachable;
//  2. Анализаторы класса SA пакета staticcheck;
//  3. Анализатор ST1013 (Should use constants for HTTP error codes, not magic numbers) пакета staticcheck;
//  4. Публичные анализаторы - errcheck, errwrap;
//  5. Собственный анализатор - osexitanalyzer, который проверяет наличие вызовов os.Exit в функции main пакета main.
//
// Использование:
//
//	$ go build -o staticlint *.go
//	$ ./staticlint ./...
package main

import (
	"github.com/fatih/errwrap/errwrap"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck/st1013"
)

func main() {
	checks := make([]*analysis.Analyzer, 0, 110)

	// Стандартные анализаторы
	checks = append(checks,
		printf.Analyzer,
		shadow.Analyzer,
		assign.Analyzer,
		defers.Analyzer,
		nilness.Analyzer,
		unreachable.Analyzer)

	// Анализаторы класса SA пакета staticcheck
	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	// Анализатор ST1013 (Should use constants for HTTP error codes, not magic numbers) пакета staticcheck
	checks = append(checks, st1013.Analyzer)

	// Публичные анализаторы
	checks = append(checks, errcheck.Analyzer)
	checks = append(checks, errwrap.Analyzer)

	// Собственный анализатор
	checks = append(checks, OsExitAnalyzer)

	multichecker.Main(
		checks...,
	)
}
