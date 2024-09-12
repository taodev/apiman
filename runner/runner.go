package runner

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/taodev/apiman/client/http"
	"github.com/taodev/apiman/logger"
	"github.com/taodev/apiman/storage"
	"gopkg.in/yaml.v3"
)

type Case []string
type Cases map[string]Case

type Runner struct {
	Cases    Cases `yaml:"cases"`
	fullpath string
	fileDir  string
	workDir  string

	ctx context.Context
}

// 获取所有用例
func (r *Runner) GetAllCases() (cases []string) {
	for k := range r.Cases {
		cases = append(cases, k)
	}
	return
}

func (r *Runner) Do(caseName string) (result *CaseResult, err error) {
	if result, err = r.run(caseName); err != nil {
		return
	}

	return
}

func (r *Runner) Load() (err error) {
	fileContent, err := os.ReadFile(r.fullpath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(fileContent, r)
	if err != nil {
		typ := reflect.TypeOf(err)
		log.Println("err: ", typ.Name())
		return
	}

	return
}

func (r *Runner) run(caseName string) (ret *CaseResult, err error) {
	caseValue, ok := r.Cases[caseName]
	if !ok {
		err = fmt.Errorf("case %s not found", caseName)
		return
	}

	ret = NewCaseResult(caseName)
	now := time.Now()
	defer func() {
		ret.Time = time.Since(now)
	}()

	sessionDB := storage.NewFromMemory()

	for _, name := range caseValue {
		select {
		case <-r.ctx.Done():
			return
		default:
			// 判断是否绝对路径
			if !filepath.IsAbs(name) {
				name = filepath.Join(r.fileDir, name)
			}

			api := new(http.ApiHttp)
			configPath := name
			if err = api.Load(r.workDir, configPath); err != nil {
				return
			}

			var result http.ApiResult

			if result, err = api.Do(sessionDB); err != nil {
				err = nil
				ret.Add(&result)
				logger.LogYaml(result)
				continue
			}

			ret.Add(&result)
			logger.LogYaml(result)
		}
	}

	return
}

func NewRunner(workDir, name string, ctx context.Context) (r *Runner) {
	r = new(Runner)
	r.ctx = ctx

	// 获取文件绝对路径
	fullpath, err := filepath.Abs(name)
	if err != nil {
		return
	}

	r.fullpath = fullpath
	r.workDir = workDir
	r.fileDir = filepath.Dir(fullpath)

	return
}
