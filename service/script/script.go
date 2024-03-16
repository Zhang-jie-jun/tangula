package script

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao"
	"github.com/Zhang-jie-jun/tangula/internal/dao/auth"
	"github.com/Zhang-jie-jun/tangula/internal/dao/script"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/Zhang-jie-jun/tangula/service"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetScripts(queryParam *view.QueryParam, user *auth.User) (int64, []map[string]interface{}, error) {
	totalNum, replicas, err := script.ScriptMgm.GetScripts(queryParam.Index,
		queryParam.Count, user.Account, queryParam.Filter)
	if err != nil {
		err = errors.New(msg.ERROR_GET_SCRIPT_FILE_INFO, msg.GetMsg(msg.ERROR_GET_SCRIPT_FILE_INFO, err.Error()))
		logrus.Error(err)
		return 0, nil, err
	}
	var result []map[string]interface{}
	for _, iter := range replicas {
		result = append(result, iter.TransformMap())
	}
	return totalNum, result, nil
}

func createCache() (string, error) {
	if _, err := os.Stat(contants.AppCfg.System.ScriptPath); err != nil {
		// not exists
		if err := os.MkdirAll(contants.AppCfg.System.ScriptPath, os.ModePerm); err != nil {
			err = errors.New(msg.ERROR_CREATE_FILE_CACHE, msg.GetMsg(msg.ERROR_CREATE_FILE_CACHE, err.Error()))
			logrus.Error(err)
			return "", err
		}
	}
	return contants.AppCfg.System.ScriptPath, nil
}

func UploadScript(context *gin.Context) (map[string]interface{}, error) {
	var err error
	var file *multipart.FileHeader
	var account string
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.UPLOAD_SCRIPT_FILE_FAILED, file.Filename, err.Error())
			service.CreateLogRecord(msg.UPLOAD_SCRIPT_FILE, file.Filename, detail, account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.UPLOAD_SCRIPT_FILE_SUCCESS, file.Filename)
			service.CreateLogRecord(msg.UPLOAD_SCRIPT_FILE, file.Filename, detail, account, contants.LOG_SUCCESS)
		}
	}()
	// 获取描述信息
	var param view.ScriptCreateParam
	if err = context.ShouldBind(&param); err != nil {
		err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error())))
		logrus.Error(err)
		return nil, err
	}
	// 获取当前登录用户名
	claims := jwt.ExtractClaims(context)
	account, ok := claims["account"].(string)
	if !ok {
		err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_INVALID_PARAMS)))
		return nil, err
	}
	// 获取文件信息
	file, err = context.FormFile("file")
	if err != nil {
		err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	// 检查文件格式
	index := strings.LastIndex(file.Filename, ".")
	logrus.Info("index:", index)
	suffix := string([]byte(file.Filename)[index:])
	if suffix != ".sh" && suffix != ".json" {
		err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE_IS_NOT_SHELL)))
		return nil, err
	}
	// 脚本文件不能大于10M
	if file.Size > 10*1024*1024 {
		err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE,
			msg.GetMsg(msg.ERROR_UPLOAD_FILE_SZIE_IS_TOO_BIG)))
		return nil, err
	}
	// 检查并创建文件缓存路径
	cacheDirName, ex := createCache()
	if ex != nil {
		err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE, ex.Error()))
		logrus.Error(err)
		return nil, err
	}
	// 创建脚本对象
	var scriptObj script.Script
	scriptObj.Name = file.Filename
	scriptObj.Desc = param.Desc
	scriptObj.Uuid = util.GenerateGuid()
	scriptObj.Label = strings.Replace(suffix, ".", "", -1)
	scriptObj.CreateUser = account
	// 使用脚本UUID命名存储文件
	fileName := scriptObj.Uuid + suffix
	dst := filepath.Join(cacheDirName, fileName)
	// 保存文件
	if err = context.SaveUploadedFile(file, dst); err != nil {
		err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	// 保存脚本对象
	scriptObj, err = script.ScriptMgm.CreateScript(scriptObj)
	if err != nil {
		err = errors.New(msg.ERROR_UPLOAD_FILE, msg.GetMsg(msg.ERROR_UPLOAD_FILE, err.Error()))
		logrus.Error(err)
		return nil, err
	}
	return scriptObj.TransformMap(), nil
}

func DownloadScript(context *gin.Context) error {
	var err error
	var scriptObj script.Script
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.DOWNLOAD_SCRIPT_FILE_FAILED, scriptObj.Name, err.Error())
			service.CreateLogRecord(msg.DOWNLOAD_SCRIPT_FILE, scriptObj.Name, detail, scriptObj.CreateUser, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.DOWNLOAD_SCRIPT_FILE_SUCCESS, scriptObj.Name)
			service.CreateLogRecord(msg.DOWNLOAD_SCRIPT_FILE, scriptObj.Name, detail, scriptObj.CreateUser, contants.LOG_SUCCESS)
		}
	}()
	var idparam view.IdParam
	if err = context.ShouldBindUri(&idparam); err != nil {
		err = errors.New(msg.ERROR_DOWNLOAD_FILE, msg.GetMsg(msg.ERROR_DOWNLOAD_FILE,
			msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error())))
		logrus.Error(err)
		return err
	}
	// 检查并创建文件缓存路径
	cacheDirName, ex := createCache()
	if ex != nil {
		err = errors.New(msg.ERROR_DOWNLOAD_FILE, msg.GetMsg(msg.ERROR_DOWNLOAD_FILE, ex.Error()))
		logrus.Error(err)
		return err
	}
	scriptObj, err = script.ScriptMgm.FindById(idparam.Id)
	if err != nil {
		err = errors.New(msg.ERROR_DOWNLOAD_FILE, msg.GetMsg(msg.ERROR_DOWNLOAD_FILE, err.Error()))
		logrus.Error(err)
		return err
	}
	downloadFile := fmt.Sprintf("%s/%s.%s", cacheDirName, scriptObj.Uuid, scriptObj.Label)
	if _, err = os.Stat(downloadFile); err != nil {
		err = errors.New(msg.ERROR_DOWNLOAD_FILE, msg.GetMsg(msg.ERROR_DOWNLOAD_FILE, err.Error()))
		logrus.Error(err)
		return err
	}
	context.Writer.Header().Add("Content-Type", "application/octet-stream")
	//强制浏览器下载
	context.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", scriptObj.Name))
	//浏览器下载或预览
	//context.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
	context.Header("Content-Transfer-Encoding", "binary")
	context.Header("Cache-Control", "no-cache")
	context.File(downloadFile)
	return nil
}

func GetScriptContent(context *gin.Context) error {
	var idparam view.IdParam
	if err := context.ShouldBindUri(&idparam); err != nil {
		err = errors.New(msg.ERROR_GET_SCRIPT_FILE_CONTENT, msg.GetMsg(msg.ERROR_GET_SCRIPT_FILE_CONTENT,
			msg.GetMsg(msg.ERROR_INVALID_PARAMS, err.Error())))
		logrus.Error(err)
		return err
	}
	// 检查并创建文件缓存路径
	cacheDirName, err := createCache()
	if err != nil {
		err = errors.New(msg.ERROR_GET_SCRIPT_FILE_CONTENT, msg.GetMsg(msg.ERROR_GET_SCRIPT_FILE_CONTENT, err.Error()))
		logrus.Error(err)
		return err
	}
	var scriptObj script.Script
	scriptObj, err = script.ScriptMgm.FindById(idparam.Id)
	if err != nil {
		err = errors.New(msg.ERROR_GET_SCRIPT_FILE_CONTENT, msg.GetMsg(msg.ERROR_GET_SCRIPT_FILE_CONTENT, err.Error()))
		logrus.Error(err)
		return err
	}
	downloadFile := fmt.Sprintf("%s/%s.%s", cacheDirName, scriptObj.Uuid, scriptObj.Label)
	if _, err = os.Stat(downloadFile); err != nil {
		err = errors.New(msg.ERROR_GET_SCRIPT_FILE_CONTENT, msg.GetMsg(msg.ERROR_GET_SCRIPT_FILE_CONTENT, err.Error()))
		logrus.Error(err)
		return err
	}
	//context.Writer.Header().Add("Content-Type", "application/octet-stream")
	context.Writer.Header().Add("Content-Type", "application/json")
	//强制浏览器下载
	//context.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", scriptObj.Name))
	//浏览器下载或预览
	context.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", scriptObj.Name))
	context.Header("Content-Transfer-Encoding", "binary")
	context.Header("Cache-Control", "no-cache")
	context.File(downloadFile)
	return nil
}

func DeleteScript(id uint, user *auth.User) error {
	var err error
	var scriptObj script.Script
	// 闭包的方式记录操作日志
	defer func() {
		if err != nil {
			detail := msg.GetOperation(msg.DELETE_SCRIPT_FILE_FAILED, scriptObj.Name, err.Error())
			service.CreateLogRecord(msg.DELETE_SCRIPT_FILE, scriptObj.Name, detail, user.Account, contants.LOG_FAILED)
		} else {
			detail := msg.GetOperation(msg.DELETE_SCRIPT_FILE_SUCCESS, scriptObj.Name)
			service.CreateLogRecord(msg.DELETE_SCRIPT_FILE, scriptObj.Name, detail, user.Account, contants.LOG_SUCCESS)
		}
	}()
	scriptObj, err = script.ScriptMgm.FindById(id)
	if err != nil {
		err = errors.New(msg.ERROR_DELETE_SCRIPT_FILE, msg.GetMsg(msg.ERROR_DELETE_SCRIPT_FILE, err.Error()))
		logrus.Error(err)
		return err
	}
	err = service.CheckResource(service.DELETE_RESOURCE, user, contants.PRIVATE, scriptObj.CreateUser)
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = script.ScriptMgm.DeleteScript(scriptObj.Id)
	if err != nil {
		err = errors.New(msg.ERROR_DELETE_SCRIPT_FILE, msg.GetMsg(msg.ERROR_DELETE_SCRIPT_FILE, err.Error()))
		logrus.Error(err)
		return err
	}
	return nil
}

func ab_access_token(serverIp string, port string, userName string, userPass string) error {
	reqUrl := fmt.Sprintf("https://%s:%s/oauths/access_token", serverIp, port)
	reqData := map[string]interface{}{
		"isEnc":          false,
		"userName":       userName,
		"userPass":       userPass,
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
		logrus.Error(reqErr)
		return reqErr
	}
	//处理返回结果
	fmt.Println(reqUrl)
	resp, respErr := client.Do(req)
	if respErr != nil {
		logrus.Error(respErr)
		return respErr
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		logrus.Error(readErr)
		return readErr
	}
	fmt.Println("resp->", string(body))
	return nil
}

func MetaCompare(MetaParam *view.MetaParam) (result map[string]interface{}, err error) {
	serverIp := MetaParam.ConsoleIp
	port := MetaParam.Port
	username := MetaParam.Username
	password := MetaParam.Password
	//登录
	ab_access_token(serverIp, port, username, password)
	res := map[string]interface{}{
		"buildId": 111,
	}

	logrus.Info(res)
	return res, err
}

func GetVmFileRecoveryReportImpl() (result map[string]interface{}, err error) {
	var res view.VmFileRecoveryResult
	err = dao.ClientDB.Raw("SELECT sequence,pass,failed,duration,startTime from case_manager.vmfile_recovery_test where id = (SELECT max(id)  FROM case_manager.vmfile_recovery_test where title=?)", "总计").Scan(&res).Error
	if err != nil {
		logrus.Error(err)
		return
	}
	/**
		for _, iter := range res {
			resMap := map[string]interface{}{
				"pass":    iter.Pass,
				"failed":  iter.Failed,
				"duration": iter.Duration,
				"sequence": iter.Sequence,
			}
			result = append(result, resMap)
		}
	    **/
	var timeLayoutStr = "2006-01-02 15:04:05" //go中的时间格式化必须是这个时间
	result = map[string]interface{}{
		"pass":      res.Pass,
		"failed":    res.Failed,
		"duration":  res.Duration,
		"sequence":  res.Sequence,
		"startTime": res.StartTime.Format(timeLayoutStr),
		"report":    fmt.Sprintf("http://10.4.108.47:8081/report/VmFileRecovery/%s.html", res.Sequence),
	}
	return
}
