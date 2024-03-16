package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/internal/vsphere"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/service/app/vmware"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/object"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
)

var keyGen func() string

func jobWhileMoni() {
	var count = 0
	for {
		if count >= 10 {
			break //如果count>=10则退出
		}
		fmt.Println("hello,world", count)
		count++
	}
}

func init() {
	keys := make(chan string)
	go func() {
		for {
			var buf [8]byte
			for i := 0; i < 8; i++ {
				buf[i] = byte(rand.Intn(26)) + byte('A')
			}
			keys <- string(buf[:])
		}
	}()
	keyGen = func() string {
		return <-keys
	}
}

func Pic(dx, dy int) [][]uint8 {
	res := make([][]uint8, dy)
	for y := range res {
		row := make([]uint8, dx)
		for x := range row {
			row[x] = uint8(x % (y + 1))
		}
		res[y] = row
	}
	return res

}

func WordCount(s string) map[string]int {
	arr := strings.Fields(s) //分割成数组
	map_s := make(map[string]int)
	fmt.Println(arr)
	for _, v := range arr {
		_, ok := map_s[v] //判断map中是否存在当前字符串为下标的值
		if ok {
			map_s[v] += 1 //存在+1
		} else {
			map_s[v] = 1 //不存在赋值为1;
		}
	}
	return map_s
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func fibonacci() func() int {
	back1, back2 := 0, 1

	return func() int {
		temp := back1
		back1, back2 = back2, (back1 + back2)
		return temp
	}
}

func ab_access_token(serverIp string) {
	reqUrl := fmt.Sprintf("https://%s:9600/oauths/access_token", serverIp)
	reqData := map[string]interface{}{
		"isEnc":          false,
		"userName":       "zzz",
		"userPass":       "passwd.com123",
		"validPwdExpire": false,
	}
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	bytesData, _ := json.Marshal(reqData)
	req, reqErr := http.NewRequest("POST", reqUrl, bytes.NewReader(bytesData))
	//增加header选项
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("referer", reqUrl)
	req.Header.Set("Content-Type", "application/json")

	if reqErr != nil {
		panic(reqErr)
	}
	//处理返回结果
	fmt.Println(reqUrl)
	resp, _ := client.Do(req)
	//body, _ := ioutil.ReadAll(resp.Body)
	res_header := resp.Header["Set-Cookie"]
	//fmt.Println(res_header)
	for _, value := range res_header {
		//fmt.Println(key,value)
		r, _ := regexp.Compile(` userId=(.*?); HttpOnly`)
		fmt.Printf(r.FindString(value))
	}
	fmt.Println()
}

func TestVmWare_GetAllVmClient() {
	vms := vmware.NewVmWare("192.168.212.52", "Administrator@vsphere.local", "passwd.com123")
	vmList, _, _ := vms.GetAllVmClient()
	for _, vm := range vmList {
		if strings.Contains(vm.Name, "细粒度-centos7") {
			fmt.Println(vm.Name, "===>", vm.Ip)
			v := object.NewVirtualMachine(vms.Client.Client, vm.VM.Reference())
			ipAddr := vsphere.IpAddr{
				Ip:       "192.168.212.111",
				Netmask:  "255.255.255.0",
				Gateway:  "192.168.212.1",
				Hostname: "tangula",
			}
			err := vms.SetIP(v, ipAddr, "linux")
			if err != nil {
				fmt.Println("设置ip报错：", err)
			}
		}
	}
}

func testFindByIp(ip string) {
	ctx := context.Background()
	vms := vmware.NewVmWare("192.168.212.52", "Administrator@vsphere.local", "passwd.com123")
	_, err := vms.FindVMByIP(ctx, vms.Client.Client, ip)
	if err != nil {
		logrus.Error(err)
	}

}

func FormatJson(data string) string {
	var out bytes.Buffer
	json.Indent(&out, []byte(data), "", "    ")
	return out.String()
}
func SerializeToJson(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}
func UnserializeFromJson(data string, result interface{}) error {
	d := json.NewDecoder(strings.NewReader(data))
	d.UseNumber()
	return d.Decode(result)
}
func FormatStruct(data interface{}) string {
	result := SerializeToJson(data)
	return FormatJson(result)
}

func main() {

	//TestVmWare_GetAllVmClient()
	//testFindByIp("192.168.212.120")
	fmt.Println(util.AesDecrypt("YXNlX2l2X2hlYWQ6X19fXzR2RWJvMm5NaVVqNHRJajIAAAAMGJYrjuqXwm7j5XtRFz0hGA=="))
	//ab_access_token("10.2.19.94")
	//var jsonContent = map[int]string{1: "zzw", 2: "qwe"}
	//logrus.Info(SerializeToJson(jsonContent))

}
