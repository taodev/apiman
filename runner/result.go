package runner

import (
	"fmt"
	"time"

	"github.com/taodev/apiman/client/http"
	"github.com/taodev/apiman/test"
)

type CaseApi struct {
	Name  string
	Time  time.Duration    `yaml:"time"`
	Tests test.TestManager `yaml:"tests,omitempty"`
	Error string           `yaml:"error,omitempty"`
}

type CaseResult struct {
	Name   string        `yaml:"name,omitempty"`
	Status string        `yaml:"status,omitempty"`
	Time   time.Duration `yaml:"time"`
	Apis   []CaseApi     `yaml:"apis,omitempty"`
}

func (cr *CaseResult) Pass() bool {
	return cr.Status == "PASS"
}

func (cr *CaseResult) Add(apiResult *http.ApiResult) {
	cr.Apis = append(cr.Apis, CaseApi{
		Name:  apiResult.Name,
		Time:  apiResult.Time,
		Tests: apiResult.Tests,
		Error: apiResult.Error,
	})

	if !apiResult.Pass() {
		cr.Status = "FAIL"
	}
}

func (cr *CaseResult) String() (v string) {
	v = fmt.Sprintf("CASE: %s %v %s\n", cr.Name, cr.Time, cr.Status)
	for _, api := range cr.Apis {
		v += fmt.Sprintf("  API: %s %v %s\n", api.Name, api.Time, api.Tests.Status)
		if api.Tests.Pass() {
			continue
		}

		if len(api.Error) > 0 {
			v += fmt.Sprintf("    ERROR: %s\n", api.Error)
		}

		for _, test := range api.Tests.Cases {
			v += fmt.Sprintf("    TEST: %s %s\n", test.Name, test.Status)
			if !test.Pass() {
				for _, e := range test.Error {
					v += fmt.Sprintf("      %s\n", e)
				}
			}
		}
	}

	return
}

func NewCaseResult(name string) *CaseResult {
	return &CaseResult{
		Name:   name,
		Status: "PASS",
	}
}
