// Package multichecker реализует множественный статический анализатор с различными компонентами.
// В пакете используется реализация golang.org/x/tools/go/analysis/multichecker
// для объединения различных компонентов в единую утилиту.
//
// Multichecker включает в себя следующие компоненты:
//
// Из пакета golang.org/x/tools/go/analysis/passes:
// assign, atomic, atomicalign, bools, composite, copylock, deepequalerrors,
// directive, errorsas, httpresponse, ifaceassert, loopclosure, lostcancel,
// nilfunc, nilness, reflectvaluecompare, shadow, shift, sigchanyzer, sortslice,
// stdmethods, stringintconv, structtag, tests, timeformat, unmarshal,
// unreachable, unsafeptr, unusedresult, unusedwrite.
//
// Дополнительная информации о конкретном компоненте, доступна на:
// https://pkg.go.dev/golang.org/x/tools/go/analysis/passes
//
// Анализаторы пакета staticheck.io:
//   - SA* (staticcheck).
//   - S* (simple).
//   - ST* (stylecheck).
//   - QF* (quickfix).
//
// Дополнительная информации о конкретном компоненте, доступна на:
// https://staticcheck.io/docs/checks/
//
// Публичные анализаторы:
//   - bodyclose находит не закрытий body HTTP-ответа блокирующий
//     повторное использование TCP-соединения (https://github.com/timakin/bodyclose).
//   - builtinprint находит использовние функцию print и println
//     для вывода отладочной информации (https://github.com/gostaticanalysis/builtinprint).
//
// Собственный анализатор:
//   - noexit проверяет использование прямого вызова os.Exit
//     в функции main пакета main.
package multichecker

import (
	"github.com/gostaticanalysis/builtinprint"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/timakin/bodyclose/passes/bodyclose"

	"github.com/a-x-a/go-metric/pkg/noexit"
)

// MulticheckerLinter структура для хранения линтеров.
type MulticheckerLinter struct {
	checkers []*analysis.Analyzer
}

// New создает новый экземпляр линтера MulticheckerLinter.
func New() MulticheckerLinter {
	// Анализаторы пакета golang.org/x/tools/go/analysis/passes.
	checkers := []*analysis.Analyzer{
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		deepequalerrors.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
	}

	// Анализаторы пакета staticcheck.io.
	for _, analyzers := range [][]*lint.Analyzer{
		staticcheck.Analyzers,
		simple.Analyzers,
		stylecheck.Analyzers,
		quickfix.Analyzers,
	} {
		for _, v := range analyzers {
			checkers = append(checkers, v.Analyzer)
		}
	}

	// Публичные анализаторы.
	checkers = append(checkers, bodyclose.Analyzer)
	checkers = append(checkers, builtinprint.Analyzer)

	// noexit анализатор.
	checkers = append(checkers, noexit.Analyzer)

	return MulticheckerLinter{checkers}
}

// Run запускает анализатор исходного кода.
func (ml MulticheckerLinter) Run() {
	multichecker.Main(ml.checkers...)
}
