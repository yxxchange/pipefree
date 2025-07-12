package test

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/yxxchange/pipefree/config"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/http"
	"github.com/yxxchange/pipefree/infra/dal"
	"github.com/yxxchange/pipefree/infra/etcd"
	http2 "net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	config.Init("../../config.yaml")
	dal.InitDB()
	etcd.InitEtcd()
	go func() {
		if err := http.LaunchServer(); err != nil {
			panic(err)
		}
	}()
	time.Sleep(1 * time.Second) // 等待服务器启动
	m.Run()
}

func TestWatchPipe(t *testing.T) {
	client := &clientTest{Endpoint: "http://localhost:" + viper.GetString("http.port")}
	go func() {
		time.Sleep(2 * time.Second)
		if err := client.RunPipe(2); err != nil {
			t.Fatalf("failed to run pipe: %v", err)
		}
	}()
	err := client.WatchPipe("default", "container", "test-pipe")
	if err != nil {
		t.Fatalf("failed to watch pipe: %v", err)
	}
}

func TestWatchPipe_TestCancel(t *testing.T) {
	wg := &sync.WaitGroup{}
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(pipeId int) {
			defer wg.Done()
			c1 := &clientTest{
				Endpoint: viper.GetString("http.endpoint"),
				Timeout:  time.Second * time.Duration(i*5),
			}
			if err := c1.WatchPipe("default", "container", "test-pipe-"+strconv.Itoa(pipeId)); err != nil {
				log.Errorf("%d failed to watch pipe: %s", pipeId, err.Error())
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("All goroutines completed\n")
	time.Sleep(1 * time.Second)
}

const (
	PatchRunPipe = "/api/v1/pipe_exec/:id"
	WatchPipe    = "/api/v1/operator/namespace/:namespace/kind/:kind"
)

type clientTest struct {
	Endpoint string
	Timeout  time.Duration
}

func (c *clientTest) RunPipe(pipeId int) error {
	path := strings.Replace(PatchRunPipe, ":id", strconv.Itoa(pipeId), 1)
	req, err := http2.NewRequest(http2.MethodPost, c.Endpoint+path, nil)
	if err != nil {
		return err
	}
	resp, err := http2.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	err = handleHttpResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to run pipe: %w", err)
	}
	return nil
}

func (c *clientTest) WatchPipe(namespace, kind, name string) error {
	path := strings.Replace(WatchPipe, ":namespace", namespace, 1)
	path = strings.Replace(path, ":kind", kind, 1)
	req, err := http2.NewRequest(http2.MethodGet, c.Endpoint+path, nil)
	if err != nil {
		return err
	}
	client := http2.Client{
		Timeout: c.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send watch request: %w", err)
	}
	err = handleHttpResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to watch pipe: %w", err)
	}
	if resp.TransferEncoding == nil || !contains(resp.TransferEncoding, "chunked") {
		return fmt.Errorf("expected chunked transfer encoding, got: %v", resp.TransferEncoding)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	for {
		event := make(map[string]interface{})
		if err := decoder.Decode(&event); err != nil {
			if err.Error() == "EOF" {
				break // 处理完所有事件后退出
			}
			return fmt.Errorf("failed to decode event: %w", err)
		}
		fmt.Printf("Received event: %v\n", event) // 打印或处理事件
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func handleHttpResponse(resp *http2.Response) error {
	if resp.StatusCode != http2.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
