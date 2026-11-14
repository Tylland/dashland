package specs

import (
	"fmt"
	"strings"
	"testing"
	"unicode"
)

type TestResult string

const (
	TestResultFailed    TestResult = "Failed"
	TestResultSucceeded TestResult = "Succeeded"
)

type StepResult struct {
	title  string
	result TestResult
	err    error
}

func NewStepResult(format string, a ...any) *StepResult {
	return &StepResult{title: fmt.Sprintf(format, a...)}
}

func (r *StepResult) Err() error {
	return r.err
}

func (r *StepResult) With(result *StepResult) *StepResult {
	if result.err != nil && r.err == nil {
		r.WithError(result.err)
	}

	return r
}

func (r *StepResult) WithResult(result TestResult) *StepResult {
	r.result = result

	return r
}

func (r *StepResult) WithError(err error) *StepResult {
	r.err = err
	r.result = TestResultFailed

	return r
}

func (r *StepResult) WithErrorMessage(msg string) *StepResult {
	r.err = fmt.Errorf("%s", msg)
	r.result = TestResultFailed

	return r
}

func (r *StepResult) WithErrorf(format string, a ...any) *StepResult {
	r.err = fmt.Errorf(format, a...)
	r.result = TestResultFailed

	return r
}

func (r *StepResult) Success() *StepResult {
	r.result = TestResultSucceeded

	return r
}

func camelAppend(str string) string {
	w := []rune(str)
	for i := len(w) - 1; i > 1; i-- {
		if unicode.IsUpper(w[i]) {
			w = append(w[:i], append([]rune{' '}, w[i:]...)...)
		}
	}
	return string(w)
}

type scenarioStepResult struct {
	title  string
	result string
}

type Scenario struct {
	t     *testing.T
	story *Story
	title string
	steps []scenarioStepResult
}

func NewScenario(t *testing.T) *Scenario {
	return &Scenario{t: t, title: camelAppend(strings.TrimPrefix(t.Name(), "Test"))}
}

func (s *Scenario) HandleStepResult(prefix string, sr *StepResult) {
	title := prefix + " " + sr.title
	result := string(sr.result)

	if sr.result == TestResultFailed {
		result += ": " + sr.err.Error()

		s.t.Fatalf("%-60s - %s", title, result)
	} else {
		s.t.Logf("%-60s - %s", title, result)
	}

	s.steps = append(s.steps, scenarioStepResult{title: title, result: result})
}

func (s *Scenario) Given(sr *StepResult) *GivenStepBlock {
	s.HandleStepResult("Given", sr)

	return NewGivenStepBlock(s)
}

func (s *Scenario) Background(sr *StepResult) *GivenStepBlock {
	s.HandleStepResult("Background", sr)

	return NewGivenStepBlock(s)
}

type BackgroundStepBlock struct {
	scenario *Scenario
}

func NewBackgroundStepBlock(scenario *Scenario) *BackgroundStepBlock {
	return &BackgroundStepBlock{scenario: scenario}
}

func (sb *BackgroundStepBlock) Given(sr *StepResult) *GivenStepBlock {
	sb.scenario.HandleStepResult("Given", sr)

	return NewGivenStepBlock(sb.scenario)
}

type GivenStepBlock struct {
	scenario *Scenario
}

func NewGivenStepBlock(scenario *Scenario) *GivenStepBlock {
	return &GivenStepBlock{scenario: scenario}
}

func (sb *GivenStepBlock) Given(sr *StepResult) *GivenStepBlock {
	sb.scenario.HandleStepResult("Given", sr)

	return sb
}

func (sb *GivenStepBlock) And(sr *StepResult) *GivenStepBlock {
	sb.scenario.HandleStepResult("  And", sr)

	return sb
}

func (sb *GivenStepBlock) When(sr *StepResult) *WhenStepBlock {
	sb.scenario.HandleStepResult("When", sr)

	return NewWhenStepBlock(sb.scenario)
}

type WhenStepBlock struct {
	scenario *Scenario
}

func NewWhenStepBlock(scenario *Scenario) *WhenStepBlock {
	return &WhenStepBlock{scenario: scenario}
}

func (sb *WhenStepBlock) Then(sr *StepResult) *ThenStepBlock {
	sb.scenario.HandleStepResult("Then", sr)

	return NewThenStepBlock(sb.scenario)
}

type ThenStepBlock struct {
	scenario *Scenario
}

func NewThenStepBlock(scenario *Scenario) *ThenStepBlock {
	return &ThenStepBlock{scenario: scenario}
}

func (sb *ThenStepBlock) And(sr *StepResult) *ThenStepBlock {
	sb.scenario.HandleStepResult("  And", sr)

	return sb
}

func (sb *ThenStepBlock) Execute() {
	fmt.Println("")
	fmt.Println("Scenario: " + sb.scenario.title)

	maxTitleWidth := 0

	for _, s := range sb.scenario.steps {
		lines := strings.Split(s.title, "\n")
		for _, line := range lines {
			if len(line) > maxTitleWidth {
				maxTitleWidth = len(line)
			}
		}
	}

	maxTitleWidth = int(maxTitleWidth/10+1) * 10

	maxWidth := maxTitleWidth + len(TestResultSucceeded)

	bar := strings.Repeat("-", maxWidth)
	fmt.Println(bar)

	for _, s := range sb.scenario.steps {
		fmt.Printf("%-*s%s\n", maxTitleWidth, s.title, s.result)
	}

	fmt.Println(bar)
}
