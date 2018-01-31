package main

import (
	"github.com/Darkera524/WinTraceTool/func/trace"
	"fmt"
	"io"
	"strings"
	"time"
	"strconv"
	"github.com/open-falcon/common/model"
	"os"
	"encoding/json"
	"bytes"
	"net/http"
)

func CronDetect(){
	for {
		iplist := GetConfig().Ip_list
		for _,ip := range iplist{
			go detect(ip)
		}
		time.Sleep(time.Duration(60) * time.Second)
	}
}

func detect(dnsserver string) error {
	command := "powershell"
	param := []string{"Test-DnsServer", "-IPAddress", dnsserver}

	cmd,reader,err := trace.ExecCommand(command, param)
	if err != nil {
		return err
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		byt := []byte(line)
		if byt[0] == 'I'{
			continue
		} else if byt[0] == '-'{
			continue
		} else if byt[0] == '1' {
			sendDetectResult(line, dnsserver)
		}

	}
	cmd.Wait()
	return nil
}

func sendDetectResult(line string, ip string){
	fmt.Print(line)
	datas := strings.Split(line, " ")
	var statistic int
	if datas[1] == "Success"{
		times := strings.Split(datas[2], ":")
		hour, err := strconv.Atoi(times[0])
		if err != nil {
			fmt.Println(err.Error())
		}

		minutes, err := strconv.Atoi(times[1])
		if err != nil {
			fmt.Println(err.Error())
		}

		second, err := strconv.Atoi(times[2])
		if err != nil {
			fmt.Println(err.Error())
		}

		statistic = second + minutes * 60 + hour * 60 * 60
	} else if datas[1] == "NoResponse"{
		statistic = -1
	} else if datas[1] == "UnknownError" {
		statistic = -2
	} else {
		statistic = -3
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err.Error())
	}
	tag := fmt.Sprintf("server=%s", ip)

	var metrics []*model.MetricValue
	singleMetric := &model.MetricValue{
		Endpoint:  hostname,
		Metric:    "dns.server.detect",
		Value:     statistic,
		Timestamp: time.Now().Unix(),
		Step:      60,
		Type:      "GAUGE",
		Tags:      tag,
	}
	metrics = append(metrics, singleMetric)
	PostToAgent(metrics)
}

func PostToAgent(metrics []*model.MetricValue) {
	if len(metrics) == 0 {
		return
	}

	contentJson, err := json.Marshal(metrics)
	if err != nil {
		fmt.Println(err.Error())
	}
	contentReader := bytes.NewReader(contentJson)
	req, err := http.NewRequest("POST", "http://127.0.0.1:1988/v1/push", contentReader)
	if err != nil {
		fmt.Println(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
}
