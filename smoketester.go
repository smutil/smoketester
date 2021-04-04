package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Global  Global
	Targets []Target `yaml:"targets"`
}

type Global struct {
	Retry         int `yaml:"retry"`
	RetryInterval int `yaml:"retryInterval"`
	StatusCode    int `yaml:"statusCode"`
	Qualitygate   int `yaml:"qualitygate"`
}

type Target struct {
	URL           string   `yaml:"url"`
	Name          string   `yaml:"name"`
	Method        string   `yaml:"method"`
	ResponseText  []string `yaml:"responseText"`
	StatusCode    int      `yaml:"statusCode"`
	Retry         int      `yaml:"retry"`
	RetryInterval int      `yaml:"retryInterval"`
	Header        []string `yaml:"header"`
	DataPath      string   `yaml:"data"`
	Username      string   `yaml:"username"`
	Password      string   `yaml:"password"`
}

type TestResult struct {
	Name   string `json:"name"`
	Result string `json:"result"`
}

var TestResults []TestResult
var Version = "develop"

const LogFormat = "\x1b[31;1m%s\x1b[0m\n"

func main() {
	var configPath string
	var version bool
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")
	flag.BoolVar(&version, "version", false, "returns the fileuploader version")
	flag.Parse()
	if version {
		fmt.Println(Version)
		return
	}
	if err := ValidateConfigPath(configPath); err != nil {
		log.Fatal(err)
	}

	config := &Config{}
	err := ReadYML(configPath, &config)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	executeTests(*config)
}

func executeTests(config Config) {
	for _, t := range config.Targets {
		if t.Retry == 0 {
			t.Retry = config.Global.Retry
		}
		if t.RetryInterval == 0 {
			t.RetryInterval = config.Global.RetryInterval
		}
		if t.StatusCode == 0 {
			t.StatusCode = config.Global.StatusCode
		}
		log.Println("Executing Test : " + t.Name)
		if t.Method == "GET" || t.Method == "POST" || t.Method == "PUT" {
			executeRequest(t)
		} else {
			log.Println("method name must be GET or POST")
			updateTestResult(t.Name, "Failed")
		}
	}

	//log.Println(TestResults)
	if config.Global.Qualitygate != 0 {
		qualitygate(config, TestResults)
	}
}

func qualitygate(config Config, TestResults []TestResult) {
	var executedTests = len(TestResults)
	var successTests int
	for _, r := range TestResults {
		if r.Result == "Success" {
			successTests = successTests + 1
		}
	}
	log.Println("executedTests : " + strconv.Itoa(executedTests) + ", success : " + strconv.Itoa(successTests) + ", qualitygate threshold : " + strconv.Itoa(config.Global.Qualitygate) + "%")
	if ((float64(successTests) / float64(executedTests)) * 100) < float64(config.Global.Qualitygate) {
		log.Println("quality gate failed")
		os.Exit(1)
	}

}

func executeRequest(t Target) {

	var isSuccess bool = true
	for i := 0; i <= t.Retry; i++ {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		var reqData io.Reader = nil
		if t.DataPath != "" {
			reqbyte, _ := ioutil.ReadFile(t.DataPath)
			reqData = bytes.NewBuffer(reqbyte)
		}

		req, err := http.NewRequest(t.Method, t.URL, reqData)
		if len(t.Header) > 0 {
			for _, header := range t.Header {
				req.Header.Set(strings.Split(header, " ")[0], strings.Split(header, " ")[1])
			}
		}

		if t.Username != "" && t.Password != "" {
			req.SetBasicAuth(t.Username, t.Password)
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Printf(LogFormat, fmt.Sprintf("error: %s", err))
			isSuccess = false
		} else {
			log.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
		}

		if isSuccess && t.StatusCode != 0 {
			log.Println("checking https status code")
			if t.StatusCode != resp.StatusCode {
				log.Println("statusCode check failed")
				isSuccess = false
			}
		}
		if isSuccess && len(t.ResponseText) > 0 {
			log.Println("checking response text")
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf(LogFormat, fmt.Sprintf("error: %s", err))
				isSuccess = false
			} else {
				for _, text := range t.ResponseText {
					if !strings.Contains(string(body), text) {
						log.Println("responseText check failed for text : " + text)
						isSuccess = false
						break
					}
				}
			}
		}
		if isSuccess == false && t.Retry > 0 && t.RetryInterval > 0 && i < t.Retry {
			log.Println("will retry after " + strconv.Itoa(t.RetryInterval) + " seconds")
			time.Sleep(time.Second * time.Duration(t.RetryInterval))
		} else {
			break
		}
	}

	if isSuccess {
		updateTestResult(t.Name, "Success")
	} else {
		updateTestResult(t.Name, "Failed")
	}

}

func updateTestResult(testName string, result string) {
	testResult := TestResult{Name: testName, Result: result}
	TestResults = append(TestResults, testResult)
}

func updateStatusAndExit(err error, testname string) {
	if err != nil {
		log.Printf(LogFormat, fmt.Sprintf("error: %s", err))
		updateTestResult(testname, "Failed")
		return
	}
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func ReadYML(configPath string, configPointer interface{}) error {
	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Init new YAML decoder
	d := yaml.NewDecoder(file)
	if err := d.Decode(configPointer); err != nil {
		return err
	}

	return nil
}

func Info(format string, args ...interface{}) {
	fmt.Printf(LogFormat, fmt.Sprintf(format, args...))
}
