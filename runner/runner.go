package runner

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

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
}

func (r *Runner) Do(workDir, name string, caseName string) (err error) {
	// 获取文件绝对路径
	fullpath, err := filepath.Abs(name)
	if err != nil {
		return
	}

	r.fullpath = fullpath
	r.workDir = workDir
	r.fileDir = filepath.Dir(fullpath)

	if err = r.load(); err != nil {
		return
	}

	if err = r.run(caseName); err != nil {
		return
	}

	return
}

func (r *Runner) load() (err error) {
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

func (r *Runner) run(caseName string) (err error) {
	caseValue, ok := r.Cases[caseName]
	if !ok {
		err = fmt.Errorf("case %s not found", caseName)
		return
	}

	sessionDB := storage.NewFromMemory()

	for _, name := range caseValue {
		// 判断是否绝对路径
		if !filepath.IsAbs(name) {
			name = filepath.Join(r.fileDir, name)
		}

		logger.Logf("==%s================================================================", name)
		api := new(http.ApiHttp)
		configPath := name
		if err = api.Load(r.workDir, configPath); err != nil {
			return
		}

		if _, err = api.Do(sessionDB); err != nil {
			return
		}

		// fmt.Println(api)
		logger.Log(api)
	}

	return
}
