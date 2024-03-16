package util

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	aes "github.com/Zhang-jie-jun/gmsm/tiny-aes"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/system"
	"github.com/Zhang-jie-jun/tangula/pkg/util/httpAuth"
	"github.com/bndr/gojenkins"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"golang.org/x/crypto/ssh"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type NcSshClient struct {
	Client *ssh.Client
}

// 生成UUID
func GenerateGuid() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// 使用MD5重试一次
		md5hash := md5.New()
		n := time.Now().String()
		md5hash.Write([]byte(n))
		uuid = hex.EncodeToString(md5hash.Sum(nil))
		return uuid
	}
	uuid = fmt.Sprintf("%x", b)
	return uuid
}

// AES加密
func AesEncrypt(cleartext string) (password string, err error) {
	secret := "abapollof97e41279b0538cef5716413"
	password, err = aes.Encrypt(secret, cleartext)
	if err != nil {
		return password, err
	}
	return password, err
}

// AES解密
func AesDecrypt(cryptograph string) (password string, err error) {
	secret := "abapollof97e41279b0538cef5716413"
	password, err = aes.Decrypt(secret, cryptograph)
	if err != nil {
		return password, err
	}
	return password, err
}

// 字符串处理
func DeleteSpace(src string) (dest string) {
	// 删除字符串中的多余空格，有多个空格时，仅保留一个空格
	if src == "" {
		return ""
	}
	// 替换tab为空格
	temp := strings.Replace(src, "  ", " ", -1)
	// 匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	dest = reg.ReplaceAllString(temp, "")
	return dest
}

func DeleteExtraSpace(src string) (dest string) {
	// 删除字符串中的多余空格，有多个空格时，仅保留一个空格
	if src == "" {
		return ""
	}
	// 替换tab为空格
	temp := strings.Replace(src, "  ", " ", -1)
	// 匹配两个及两个以上空格的正则表达式
	reg := regexp.MustCompile("\\s{2,}")
	dest = reg.ReplaceAllString(temp, "")
	return dest
}

func MapToJson(outputData map[string]interface{}) (jsonStr string, err error) {
	temp, err := json.Marshal(outputData)
	if err != nil {
		return jsonStr, err
	}
	jsonStr = string(temp)
	return jsonStr, err
}

func JsonToMap(jsonStr string) (outputData map[string]interface{}, err error) {
	err = json.Unmarshal([]byte(jsonStr), &outputData)
	if err != nil {
		return outputData, err
	}
	return outputData, err
}

func MapInterfaceToMapString(inputData map[string]interface{}) (outputData map[string]string, err error) {
	outputData = make(map[string]string, 0)
	for key, value := range inputData {
		switch value.(type) {
		case string:
			outputData[key] = value.(string)
		case int:
			tmp := value.(int)
			outputData[key] = cast.ToString(tmp)
		case int64:
			tmp := value.(int64)
			outputData[key] = cast.ToString(tmp)
		case uint:
			tmp := value.(uint)
			outputData[key] = cast.ToString(tmp)
		case uint64:
			tmp := value.(uint64)
			outputData[key] = cast.ToString(tmp)
		case float32:
			tmp := value.(float32)
			outputData[key] = cast.ToString(tmp)
		case float64:
			tmp := value.(float64)
			outputData[key] = cast.ToString(tmp)
		case bool:
			tmp := value.(bool)
			outputData[key] = cast.ToString(tmp)
		case byte:
			tmp := value.(byte)
			outputData[key] = cast.ToString(tmp)
		case []int:
			tmp := value.([]int)
			tmp1, _ := json.Marshal(tmp)
			outputData[key] = string(tmp1)
		case []string:
			tmp := value.([]string)
			tmp1, _ := json.Marshal(tmp)
			outputData[key] = string(tmp1)
		case map[string]string:
			tmp := value.(map[string]string)
			tmp1, _ := json.Marshal(tmp)
			outputData[key] = string(tmp1)
		case map[string]interface{}:
			tmp := value.(map[string]interface{})
			tmp1, _ := json.Marshal(tmp)
			outputData[key] = string(tmp1)
		default:
			err = errors.New("nonsupport type!")
			return outputData, err
		}
	}
	return outputData, err
}

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// 字符串去除空格
func Strips(s_ string, chars_ string) string {
	s, chars := []rune(s_), []rune(chars_)
	length := len(s)
	max := len(s) - 1
	l, r := true, true //标记当左端或者右端找到正常字符后就停止继续寻找
	start, end := 0, max
	tmpEnd := 0
	charset := make(map[rune]bool) //创建字符集，也就是唯一的字符，方便后面判断是否存在
	for i := 0; i < len(chars); i++ {
		charset[chars[i]] = true
	}
	for i := 0; i < length; i++ {
		if _, exist := charset[s[i]]; l && !exist {
			start = i
			l = false
		}
		tmpEnd = max - i
		if _, exist := charset[s[tmpEnd]]; r && !exist {
			end = tmpEnd
			r = false
		}
		if !l && !r {
			break
		}
	}
	if l && r { // 如果左端和右端都没找到正常字符，那么表示该字符串没有正常字符
		return ""
	}
	return string(s[start : end+1])
}

/********************/
//处理中文乱码的问题
type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(byte []byte, charset Charset) string {

	var str string
	switch charset {
	case GB18030:
		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}

	return str
}

/********************/

type AuthRes struct {
	StatusCode int
	StatusMap  map[string]interface{}
}

func CasRequest(userName string, password string, method string, reqUrl string, reqData map[string]interface{}) (response *AuthRes, err error) {
	logrus.Info(fmt.Sprintf("【%s】请求参数:%s", reqUrl, SerializeToJson(reqData)))
	var res AuthRes
	t := httpAuth.NewTransport(userName, password)
	bytesData, _ := json.Marshal(reqData)
	req, reqErr := http.NewRequest(method, reqUrl, bytes.NewReader(bytesData))
	if reqErr != nil {
		logrus.Error(reqErr)
		return &res, reqErr
	}
	resp, err := t.RoundTrip(req)
	if err != nil {
		logrus.Error(err)
		return &res, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Info("StatusCode=====>", resp.StatusCode)
	logrus.Info("ResBody========>", string(body))
	res.StatusCode = resp.StatusCode
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)
	res.StatusMap = resMap
	logrus.Info(fmt.Sprintf("【%s】请求返回:%s", reqUrl, SerializeToJson(res)))
	if resp.StatusCode != 200 {
		msg_str := cast.ToString(resp.Header["Error-Message"])
		w := []byte(msg_str)
		garbledStr := ConvertByte2String(w, GB18030)
		err = errors.New(fmt.Sprintf("【%s】请求返回异常：%s", reqUrl, garbledStr))
		return &res, err
	}
	defer resp.Body.Close()
	return &res, err
}

type Casjson struct {
	Id          string                   `json:"id"`
	Name        string                   `json:"name"`
	Title       string                   `json:"title"`
	DiskList    []map[string]interface{} `json:"disk_list"`
	NetworkList []map[string]interface{} `json:"network_list"`
	OsVersion   string                   `json:"osVersion"`
	System      string                   `json:"system"`
}

func ReadCasJson(fileName string) ([]Casjson, error) {
	u := []Casjson{}
	bytes, readErr := ioutil.ReadFile(fileName)
	if readErr != nil {
		logrus.Error("读取json文件失败:", readErr)
		return u, readErr
	}

	marshallErr := json.Unmarshal(bytes, &u)
	if marshallErr != nil {
		logrus.Error("解析json数据失败:", marshallErr)
		return u, marshallErr
	}
	res := fmt.Sprintf("%+v\n", u)
	logrus.Info(res)
	return u, nil
}

type Metajson struct {
	Version      uint                     `json:"version"`
	PlatformType uint                     `json:"platformType"`
	CpuSlots     uint                     `json:"cpuSlots"`
	CpuCores     uint                     `json:"cpuCores"`
	MemorySize   uint                     `json:"memorySize"`
	VmName       string                   `json:"vmName"`
	Disk         []map[string]interface{} `json:"disk"`
	Nic          []map[string]interface{} `json:"nic"`
}

func ReadMetaJson(fileName string) (Metajson, error) {
	u := Metajson{}
	bytes, readErr := ioutil.ReadFile(fileName)
	if readErr != nil {
		logrus.Error("读取json文件失败:", readErr)
		return u, readErr
	}

	marshallErr := json.Unmarshal(bytes, &u)
	if marshallErr != nil {
		logrus.Error("解析json数据失败:", marshallErr)
		return u, marshallErr
	}
	res := fmt.Sprintf("%+v\n", u)
	logrus.Info(res)
	return u, nil
}

func IsFileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 触发jenkins任务
func BuildJenkins(context *gin.Context, jobName string, paramMap map[string]string) (buildId int64, err error) {
	jenkins := gojenkins.CreateJenkins(nil, contants.AppCfg.JENKINS.Url, contants.AppCfg.JENKINS.Username, contants.AppCfg.JENKINS.Password)
	_, initErr := jenkins.Init(context)
	if initErr != nil {
		logrus.Error(initErr)
		return buildId, initErr
	}
	logrus.Info("jenkins is ok")

	logrus.Info("开始执行jenkins job:", jobName)
	logrus.Info(paramMap)
	_, buildErr := jenkins.BuildJob(context, jobName, paramMap)
	if buildErr != nil {
		logrus.Error(buildErr)
		return buildId, buildErr
	}

	job, getErr := jenkins.GetJob(context, jobName)
	if getErr != nil {
		logrus.Error(getErr)
	}

	lastBuild, lastErr := job.GetLastBuild(context)
	if lastErr != nil {
		logrus.Error(lastErr)
	}
	logrus.Info("buildID:", lastBuild.Info())

	return lastBuild.GetBuildNumber(), nil
}

// map输出格式化
func FormatJson(data interface{}) string {
	var in []byte
	switch temp := data.(type) {
	case []byte:
		in = temp
	case string:
		in = []byte(temp)
	default:
		in, _ = json.Marshal(data)

	}
	var out bytes.Buffer
	err := json.Indent(&out, in, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	_, _ = out.WriteTo(os.Stdout)
	return out.String()
}

func SerializeToJson(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}

// 获取磁盘信息
func GetMountDiskInfo(diskPath string) ([]string, error) {
	// getDiskInfoCmd := fmt.Sprintf("ls -l %s |grep 'raw' | awk '/.raw/&&!/.vmdk/{print$5,$9}'", diskPath)
	getDiskInfoCmd := fmt.Sprintf("ls -l  %s | awk '{print$5,$9}'", diskPath)
	result, rbdErr := system.SysManage.RunCommand(getDiskInfoCmd)
	if rbdErr != nil {
		logrus.Error(rbdErr)
		return nil, rbdErr
	}
	return result, nil

}
