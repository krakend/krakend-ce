package tests

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"
)

var (
	defaultBinPath   *string = flag.String("krakend_bin_path", ".././krakend", "The default path to the krakend bin")
	defaultCfgPath   *string = flag.String("krakend_config_path", "fixtures/krakend.json", "The default path to the krakend config")
	defaultSpecsPath *string = flag.String("krakend_specs_path", "./fixtures/specs", "The default path to the specs folder")
)

func TestIntegration(t *testing.T) {
	cmd := exec.Command(*defaultBinPath, "run", "-d", "-c", *defaultCfgPath)

	if len(cmd.Env) == 0 {
		cmd.Env = []string{}
	}
	cmd.Env = append(cmd.Env, "USAGE_DISABLE=1")

	if err := cmd.Start(); err != nil {
		t.Error(err)
		return
	}
	defer cmd.Process.Kill()

	go func() { fmt.Println(cmd.Wait()) }()

	tcs, err := testCases()
	if err != nil {
		t.Error(err)
		return
	}

	backend := newMockBackend()
	defer backend.Close()

	go func() {
		if err := backend.ListenAndServe(); err != nil {
			log.Printf("backend closed: %v", err)
		}
	}()

	select {
	case <-time.After(1500 * time.Millisecond):
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			req, err := newRequest(tc.In)
			if err != nil {
				t.Error(err)
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil && err.Error() != tc.Err {
				t.Error(err)
				return
			}

			if err != nil {
				return
			}

			if resp.StatusCode != tc.Out.StatusCode {
				t.Errorf("unexpected status code. have: %d, want: %d", resp.StatusCode, tc.Out.StatusCode)
			}

			for k, v := range tc.Out.Header {
				if h := resp.Header.Get(k); h != v {
					t.Errorf("unexpected value for header %s. have: %s, want: %s", k, h, v)
				}
			}

			body := ""

			if resp.Body != nil {
				b, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
					return
				}
				resp.Body.Close()
				body = string(b)
			}

			if tc.Out.Body != body {
				t.Errorf("unexpected body. have: %s\n\twant: %s", body, tc.Out.Body)
			}
		})
	}
}

func testCases() ([]TestCase, error) {
	tcs := []TestCase{}
	content, err := readSpecs()
	if err != nil {
		return tcs, err
	}

	for name, c := range content {
		tc, err := parseTestCase(name, c)
		if err != nil {
			return tcs, err
		}
		tcs = append(tcs, tc)
	}

	return tcs, nil
}

func parseTestCase(name string, in []byte) (TestCase, error) {
	tc := TestCase{}
	if err := json.Unmarshal(in, &tc); err != nil {
		return tc, err
	}
	tc.Name = name

	return tc, nil
}

func newRequest(in In) (*http.Request, error) {
	var body io.Reader
	if in.Body != "" {
		body = bytes.NewBufferString(in.Body)
	}
	req, err := http.NewRequest(in.Method, in.URL, body)
	if err != nil {
		return nil, err
	}
	for k, v := range in.Header {
		req.Header.Add(k, v)
	}
	return req, nil
}

func readSpecs() (map[string][]byte, error) {
	data := map[string][]byte{}
	files, err := ioutil.ReadDir(*defaultSpecsPath)
	if err != nil {
		return data, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		content, err := ioutil.ReadFile(path.Join(*defaultSpecsPath, file.Name()))
		if err != nil {
			return data, err
		}
		data[file.Name()[:len(file.Name())-5]] = content
	}
	return data, nil
}

func newMockBackend() http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/param_forwarding/", echoEndpoint)
	mux.HandleFunc("/delayed/", func(rw http.ResponseWriter, r *http.Request) {
		<-time.After(200 * time.Millisecond)
		echoEndpoint(rw, r)
	})

	return http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
}

func echoEndpoint(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	if ip := net.ParseIP(r.Header.Get("X-Forwarded-For")); ip == nil || !ip.IsLoopback() {
		http.Error(rw, "invalid X-Forwarded-For", 400)
		return
	}
	r.Header.Del("X-Forwarded-For")
	json.NewEncoder(rw).Encode(map[string]interface{}{
		"path":    r.URL.Path,
		"query":   r.URL.Query(),
		"headers": r.Header,
		"foo":     42,
	})
}

type TestCase struct {
	Name string `json:"name"`
	Err  string `json:"error"`
	In   In     `json:"in"`
	Out  Out    `json:"out"`
}

type In struct {
	URL    string            `json:"url"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`
}

type Out struct {
	StatusCode int               `json:"status_code"`
	Body       string            `json:"body"`
	Header     map[string]string `json:"header"`
}
