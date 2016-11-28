package owl

import (
	"github.com/pkg/errors"
	runner_python "github.com/yi-jiayu/esd-tutor-runner-python/lib"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"log"
)

type Job struct {
	Language  string     `json:"language"`
	Source    string     `json:"source_code"`
	TestCases []TestCase `json:"test_cases"`
}

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Types  string `json:"types"`
}

type Results struct {
	NumTests int `json:"num_tests"`
	Passed   int `json:"passed"`
	Failed   int `json:"failed"`
	Errored  int `json:"errored"`
}

func RunJob(job Job) (results Results, err error) {
	// check if language is supported
	// currenly only python is supported
	if job.Language != "python" {
		return results, errors.New("unsupported language")
	}

	// create a temp directory to store all the files we are going to create for easy removal
	testDir, err := ioutil.TempDir("", "job")
	if err != nil {
		return results, err
	}

	defer func() {
		if r := recover(); r != nil {
			// testDir needs to be cleaned up
			log.Println(r)
			err2 := os.RemoveAll(testDir)
			if err2 != nil {
				// double fault??
				err = err2
			}

			// while debugging don't clean up the temp files, instead we want to inspect the directory
			//log.Printf("temp files are stored at: %s", testDir)
		}
	}()

	// write source to temp file
	sourceFile, err := ioutil.TempFile(testDir, "source")
	if err != nil {
		panic(err)
	}

	_, err = sourceFile.WriteString(job.Source)
	if err != nil {
		panic(err)
	}

	// set read permission for temp file to universe
	err = os.Chmod(sourceFile.Name(), 444)
	if err != nil {
		panic(err)
	}

	numTests := len(job.TestCases)
	passed, failed, errored := 0, 0, 0

	for _, testCase := range job.TestCases {
		// create temporary file for program output
		outFile, err := ioutil.TempFile(testDir, "case")
		if err != nil {
			panic(err)
		}

		// run test
		runner := runner_python.Runner{
			SourceFile: sourceFile,
			In:         strings.NewReader(testCase.Input),
			Out:        outFile,
			Timeout:    10 * time.Second,
			Image:      "python:3.5.2-slim",
		}

		err = runner.Run()
		if err != nil {
			errored++
			continue
		}

		// check output
		types := []int{}
		for _, t := range strings.Split(testCase.Types, " ") {
			switch t {
			case "string":
				types = append(types, STR)
			case "int":
				types = append(types, INT)
			case "float":
				types = append(types, FP)
			default:
				types = append(types, STR)
			}
		}

		// looks like we need to seek back to the beginning of outFile
		_, err = outFile.Seek(0, os.SEEK_SET)
		if err != nil {
			return results, errors.Wrap(err, "error while seeking outFile")
		}

		result, err := Check(strings.NewReader(testCase.Output), outFile, types)
		if err != nil {
			return results, err
		}

		if result {
			passed++
		} else {
			failed++
		}
	}

	// clean up temp files
	err = os.RemoveAll(testDir)
	if err != nil {
		return results, err
	}
	// while debugging don't clean up the temp files, instead we want to inspect the directory
	//log.Printf("temp files are stored at: %s", testDir)

	results = Results{
		NumTests: numTests,
		Passed:   passed,
		Failed:   failed,
		Errored:  errored,
	}

	return results, nil
}
