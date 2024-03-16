package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/service/app/fusioncompute"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"io/ioutil"
	"net/http"
)

func get_version(serverIp string) (response *fusioncompute.FcResult, err error) {
	var fcRes fusioncompute.FcResult
	reqUrl := fmt.Sprintf("https://%s:7443/service/versions", serverIp)
	req, reqErr := http.NewRequest("GET", reqUrl, nil)
	//增加header选项
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Accept", "application/json;charset=UTF-8")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", serverIp))

	if reqErr != nil {
		return &fcRes, reqErr
	}
	//处理返回结果
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)
	logrus.Info(fmt.Sprintf("【%s】=>StatusCode:%d", reqUrl, resp.StatusCode))
	logrus.Info(fmt.Sprintf("【%s】=>resBody:%s", reqUrl, util.SerializeToJson(resMap)))
	fcRes.StatusCode = resp.StatusCode
	fcRes.ResBody = resMap
	return &fcRes, err
}

func login(serverIp string) (token string, err error) {
	reqUrl := fmt.Sprintf("https://%s:7443/service/session", serverIp)
	req, reqErr := http.NewRequest("POST", reqUrl, nil)
	//增加header选项
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.0")
	req.Header.Set("Content-Type", "appilcation/json;charset=UTF-8")
	req.Header.Set("X-Auth-User", "Duke")
	req.Header.Set("X-Auth-Key", "passwd.com123")
	req.Header.Set("X-Auth-UserType", cast.ToString(2))
	req.Header.Set("X-ENCRIPT-ALGORITHM", cast.ToString(1))

	if reqErr != nil {
		return "", reqErr
	}
	//处理返回结果
	fmt.Println(reqUrl)

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)
	logrus.Info(fmt.Sprintf("【%s】=>StatusCode:%d", reqUrl, resp.StatusCode))
	logrus.Info(fmt.Sprintf("【%s】=>resBody:%s", reqUrl, util.SerializeToJson(resMap)))
	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("【%s】请求返回异常：%s", reqUrl, util.SerializeToJson(resMap)))
		return "", err
	} else {
		return resp.Header["X-Auth-Token"][0], nil
	}

}

func get_site(serverIp string) {
	token, err := login(serverIp)
	if err != nil {
		logrus.Error(err)
	}
	reqUrl := fmt.Sprintf("https://%s:7443/service/sites", serverIp)
	req, _ := http.NewRequest("GET", reqUrl, nil)

	//处理返回结果
	fmt.Println(reqUrl)
	//增加header选项
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.1")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", serverIp))
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)

	fmt.Println(fmt.Sprintf("StatusCode:%d", resp.StatusCode))
	fmt.Println(util.SerializeToJson(resMap))
}

func addDataresource(serverIp string) {
	token, err := login(serverIp)
	if err != nil {
		logrus.Error(err)
	}
	reqUrl := fmt.Sprintf("https://%s:7443/service/sites/13531140/storageresources", serverIp)
	data_channel_list := []map[string]interface{}{}
	data_channel_map := map[string]interface{}{"ip": "10.4.117.167"}
	data_channel_list = append(data_channel_list, data_channel_map)
	var host_urn_list []string
	host_urn_list = append(host_urn_list, "urn:sites:13531140:hosts:105")
	reqParam := map[string]interface{}{
		"name":        "tangula1",
		"dataChannel": data_channel_list,
		"hostUrn":     host_urn_list,
		"storageType": "NAS",
		"vender":      "OTHER",
		"deviceType":  "OTHER",
		"autoscan":    1,
	}
	bytesData, byteErr := json.Marshal(reqParam)
	if byteErr != nil {
		fmt.Println(byteErr)
	}
	req, reqErr := http.NewRequest("POST", reqUrl, bytes.NewReader(bytesData))
	if reqErr != nil {
		fmt.Println(reqErr)
	}
	//处理返回结果
	fmt.Println(reqUrl)
	//增加header选项
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.0")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", serverIp))
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)

	fmt.Println(fmt.Sprintf("StatusCode:%d", resp.StatusCode))
	fmt.Println(util.SerializeToJson(resMap))

}

func get_storageresources(serverIp string) {
	token, err := login(serverIp)
	if err != nil {
		logrus.Error(err)
	}
	reqUrl := fmt.Sprintf("https://%s:7443/service/sites/13531140/storageresources?limit=5&offset=0&name=tangula", serverIp)
	req, reqErr := http.NewRequest("GET", reqUrl, nil)
	if reqErr != nil {
		fmt.Println(reqErr)
	}
	//处理返回结果
	fmt.Println(reqUrl)
	//增加header选项
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.0")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", serverIp))
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)

	fmt.Println(fmt.Sprintf("StatusCode:%d", resp.StatusCode))
	fmt.Println(util.SerializeToJson(resMap))
}

func connectHost(serverIp string, storageId uint) {
	token, err := login(serverIp)
	if err != nil {
		logrus.Error(err)
	}
	reqUrl := fmt.Sprintf("https://%s:7443/service/sites/13531140/storageresources/%d/action/connect?resID=%d", serverIp, storageId, storageId)
	reqParam := map[string]interface{}{
		"hostUrn": "urn:sites:13531140:hosts:105",
	}
	bytesData, byteErr := json.Marshal(reqParam)
	if byteErr != nil {
		fmt.Println(byteErr)
	}
	req, reqErr := http.NewRequest("POST", reqUrl, bytes.NewReader(bytesData))
	if reqErr != nil {
		fmt.Println(reqErr)
	}
	//处理返回结果
	fmt.Println(reqUrl)
	//增加header选项
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.0")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", serverIp))
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)

	fmt.Println(fmt.Sprintf("StatusCode:%d", resp.StatusCode))
	fmt.Println(util.SerializeToJson(resMap))

}

func querystorageresourcehosts(serverIp string, storageId uint, hostIp string) {
	token, err := login(serverIp)
	if err != nil {
		logrus.Error(err)
	}
	reqUrl := fmt.Sprintf("https://%s:7443/service/sites/13531140/storageresources/%d/querystorageresourcehosts?limit=5&offset=0&status=1&ip=%s", serverIp, storageId, hostIp)
	req, reqErr := http.NewRequest("GET", reqUrl, nil)
	if reqErr != nil {
		fmt.Println(reqErr)
	}
	//处理返回结果
	fmt.Println(reqUrl)

	//增加header选项
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.0")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", serverIp))
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)

	fmt.Println(fmt.Sprintf("StatusCode:%d", resp.StatusCode))
	fmt.Println(util.SerializeToJson(resMap))
}

func refresh(serverIp string, hostUrn string) {
	token, err := login(serverIp)
	if err != nil {
		logrus.Error(err)
	}
	reqUrl := fmt.Sprintf("https://%s:7443/service/sites/13531140/storageunits/action/refresh", serverIp)
	reqParam := map[string]interface{}{
		"hostUrn": hostUrn,
	}
	bytesData, byteErr := json.Marshal(reqParam)
	if byteErr != nil {
		fmt.Println(byteErr)
	}
	req, reqErr := http.NewRequest("POST", reqUrl, bytes.NewReader(bytesData))
	if reqErr != nil {
		fmt.Println(reqErr)
	}
	//处理返回结果
	fmt.Println(reqUrl)
	//增加header选项
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.0")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", serverIp))
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)

	fmt.Println(fmt.Sprintf("StatusCode:%d", resp.StatusCode))
	fmt.Println(util.SerializeToJson(resMap))

	//等待刷新完成
	if resp.StatusCode != 200 {
		if resMap["errorCode"] != "10410023" {
			logrus.Error(util.SerializeToJson(resMap))
			return
		}
	}
	var count = 0
	for {
		logrus.Info(fmt.Sprintf("第%d次查询存储设备……", count+1))
		if count > 10 {
			logrus.Error("100秒后仍未识别到存储设备")
			break
		}
		queryFlag := false
		queryRes, queryErr := queryallstorageunit(serverIp)
		if queryErr != nil {
			logrus.Error(fmt.Sprintf("查询存储设备报错：%s", queryErr.Error()))
		} else {
			suList := cast.ToSlice(queryRes["suList"])
			for _, iter := range suList {
				if cast.ToStringMap(iter)["name"] == "10.4.117.167:/tangula/mnt/92f7efae45a97e644400c209be519461/68466067b6cbf2d3d3f0170f8cce79e7" {
					logrus.Error(fmt.Sprintf("已识别到存储设备：%s", cast.ToStringMap(iter)["name"]))
					queryFlag = true
					break
				}
			}
		}
		if queryFlag {
			break
		}
		count++
	}

}

func queryallstorageunit(serverIp string) (result map[string]interface{}, err error) {
	token, err := login(serverIp)
	if err != nil {
		logrus.Error(err)
	}
	reqUrl := fmt.Sprintf("https://%s:7443/service/sites/13531140/storageunits/queryallstorageunit?limit=99&offset=0&useState=all&deviceType=0&type=NAS&name=tangula", serverIp)
	req, reqErr := http.NewRequest("GET", reqUrl, nil)
	if reqErr != nil {
		fmt.Println(reqErr)
	}
	//处理返回结果
	fmt.Println(reqUrl)

	//增加header选项
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.0")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", serverIp))
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)

	fmt.Println(fmt.Sprintf("StatusCode:%d", resp.StatusCode))
	fmt.Println(util.SerializeToJson(resMap))
	return resMap, err
}

func addDatastore(serverIp string, name string, hostUrn string, storageUnitUrn string) {
	token, err := login(serverIp)
	if err != nil {
		logrus.Error(err)
	}
	reqUrl := fmt.Sprintf("https://%s:7443/service/sites/13531140/datastores", serverIp)
	reqParam := map[string]interface{}{
		"hostUrn":        hostUrn,
		"name":           name,
		"storageUnitUrn": storageUnitUrn,
		"useType":        1,
	}
	bytesData, byteErr := json.Marshal(reqParam)
	if byteErr != nil {
		fmt.Println(byteErr)
	}
	req, reqErr := http.NewRequest("POST", reqUrl, bytes.NewReader(bytesData))
	if reqErr != nil {
		fmt.Println(reqErr)
	}
	//处理返回结果
	fmt.Println(reqUrl)
	//增加header选项
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.0")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", serverIp))
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)

	fmt.Println(fmt.Sprintf("StatusCode:%d", resp.StatusCode))
	fmt.Println(util.SerializeToJson(resMap))

}

func delDatastore(serverIp string, hostUrn string, datastoreId uint) {
	token, err := login(serverIp)
	if err != nil {
		logrus.Error(err)
	}
	reqUrl := fmt.Sprintf("https://%s:7443/service/sites/13531140/datastores/%d?hostUrn=%s&isForce=false&dataStoreID=%d", serverIp, datastoreId, hostUrn, datastoreId)

	req, reqErr := http.NewRequest("DELETE", reqUrl, nil)
	if reqErr != nil {
		fmt.Println(reqErr)
	}
	//处理返回结果
	fmt.Println(reqUrl)
	//增加header选项
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.0")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", serverIp))
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)

	fmt.Println(fmt.Sprintf("StatusCode:%d", resp.StatusCode))
	fmt.Println(util.SerializeToJson(resMap))

}

func queryDatastore(serverIp string) (result map[string]interface{}, err error) {
	token, err := login(serverIp)
	if err != nil {
		logrus.Error(err)
	}
	reqUrl := fmt.Sprintf("https://%s:7443/service/sites/13531140/datastores/ex?page=1&offset=0&limit=5&name=testtttt", serverIp)
	req, reqErr := http.NewRequest("GET", reqUrl, nil)
	if reqErr != nil {
		fmt.Println(reqErr)
	}
	//处理返回结果
	fmt.Println(reqUrl)

	//增加header选项
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.0")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", serverIp))
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)

	fmt.Println(fmt.Sprintf("StatusCode:%d", resp.StatusCode))
	fmt.Println(util.SerializeToJson(resMap))
	return resMap, err
}

func getTasks(serverIp string) (result map[string]interface{}, err error) {
	token, err := login(serverIp)
	if err != nil {
		logrus.Error(err)
	}
	reqUrl := fmt.Sprintf("https://%s:7443/service/sites/3D9907E9/tasks?page=1&limit=10&offset=0&type=RefreshStorageUnitTask", serverIp)
	req, reqErr := http.NewRequest("GET", reqUrl, nil)
	if reqErr != nil {
		fmt.Println(reqErr)
	}
	//处理返回结果
	fmt.Println(reqUrl)

	//增加header选项
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.0")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", serverIp))
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)

	fmt.Println(fmt.Sprintf("StatusCode:%d", resp.StatusCode))
	fmt.Println(util.SerializeToJson(resMap))
	return resMap, err
}

func crtVm(serverIp string) {
	token, err := login(serverIp)
	if err != nil {
		logrus.Error(err)
	}
	reqUrl := fmt.Sprintf("https://%s:7443/service/sites/13531140/vms?siteID=13531140", serverIp)

	diskList := []map[string]interface{}{
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    1,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    2,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    3,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    4,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    5,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    6,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    7,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    8,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    9,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    10,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    11,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    12,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    13,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    14,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    15,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    16,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    17,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    18,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    19,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
		{
			"datastoreUrn":   "urn:sites:13531140:datastores:91",
			"quantityGB":     1,
			"sequenceNum":    20,
			"type":           "normal",
			"indepDisk":      false,
			"persistentDisk": true,
			"isThin":         false,
			"volType":        0,
			"pciType":        "VIRTIO",
			"bootOrder":      -1,
			"encrypted":      "0",
		},
	}
	nicList := []map[string]interface{}{
		{
			"enableSecurityGroup": false,
			"portGroupUrn":        "urn:sites:13531140:dvswitchs:1:portgroups:1",
			"sequenceNum":         0,
			"virtIo":              1,
			"bootOrder":           -1,
			"nicConfig": map[string]interface{}{
				"vringbuf": 256,
				"queues":   1,
			},
		},
	}
	reqParam := map[string]interface{}{
		"name":          "异构-磁盘很多",
		"description":   "",
		"location":      "urn:sites:13531140:hosts:105",
		"parentObjUrn":  "urn:sites:13531140",
		"enableImc":     false,
		"isBindingHost": false,
		"osOptions": map[string]interface{}{
			"osType":    "Linux",
			"osVersion": 462,
		},
		"vmConfig": map[string]interface{}{
			"cpu": map[string]interface{}{
				"cpuHotPlug":      0,
				"cpuPolicy":       "shared",
				"cpuThreadPolicy": "prefer",
				"weight":          500,
				"quantity":        1,
				"limit":           0,
				"reservation":     0,
				"coresPerSocket":  1,
			},
			"memory": map[string]interface{}{
				"quantityMB":  2048,
				"hugePage":    "4K",
				"limit":       0,
				"reservation": 0,
				"weight":      20480,
			},
			"disks": diskList,
			"nics":  nicList,
			"graphicsCard": map[string]interface{}{
				"type": "cirrus",
				"size": 4,
			},
			"properties": map[string]interface{}{
				"antivirusMode":      "",
				"isAutoAdjustNuma":   false,
				"bootFirmware":       "BIOS",
				"bootFirmwareTime":   0,
				"bootOption":         "disk",
				"clockMode":          "freeClock",
				"vmVncKeymapSetting": 7,
				"isHpet":             false,
				"isEnableHa":         false,
				"evsAffinity":        false,
				"secureVmType":       "",
				"realtime":           false,
				"isEnableMemVol":     false,
				"isEnableFt":         false,
				"isAutoUpgrade":      true,
				"emulatorResType":    nil,
				"dpiVmType":          "",
				"cdRomBootOrder":     -1,
				"attachType":         false,
				"enableWatchDog":     false,
			},
			"gpuGroups": []map[string]interface{}{},
		},
		"cpuVendor":         "Intel",
		"autoBoot":          false,
		"isSrcTemplate":     false,
		"isEnableIntegrity": false,
	}
	bytesData, byteErr := json.Marshal(reqParam)
	if byteErr != nil {
		fmt.Println(byteErr)
	}
	req, reqErr := http.NewRequest("POST", reqUrl, bytes.NewReader(bytesData))
	if reqErr != nil {
		fmt.Println(reqErr)
	}
	//处理返回结果
	fmt.Println(reqUrl)
	//增加header选项
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.0")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", serverIp))
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)

	fmt.Println(fmt.Sprintf("StatusCode:%d", resp.StatusCode))
	fmt.Println(util.SerializeToJson(resMap))
}

func getVmInfo(serverIp string) (result map[string]interface{}, err error) {
	token, err := login(serverIp)
	if err != nil {
		logrus.Error(err)
	}
	reqUrl := fmt.Sprintf("https://%s:7443/service/sites/13531140/vms/i-0000012D?siteID=13531140&vmID=i-0000012D", serverIp)
	req, reqErr := http.NewRequest("GET", reqUrl, nil)
	if reqErr != nil {
		fmt.Println(reqErr)
	}
	//处理返回结果
	fmt.Println(reqUrl)

	//增加header选项
	req.Header.Set("Accept", "application/json;charset=UTF-8;version=8.0")
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", fmt.Sprintf("%s:7443", serverIp))
	req.Header.Set("X-Auth-Token", token)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resMap := make(map[string]interface{})
	json.Unmarshal(body, &resMap)

	fmt.Println(fmt.Sprintf("StatusCode:%d", resp.StatusCode))
	fmt.Println(util.SerializeToJson(resMap))
	return resMap, err
}
func main() {
	//res,err:=get_version("10.4.104.182")
	//if err !=nil{
	//	fmt.Println(err)
	//}

	//_,err:=login("10.4.104.182")
	//if err !=nil{
	//	fmt.Println(err)
	//}
	//get_site("10.4.104.182")
	//addDataresource("10.4.104.182")
	//get_storageresources("10.4.104.182")
	//connectHost("10.4.104.182", 533)
	//querystorageresourcehosts("10.4.104.182", 533,"10.4.117.95")
	//refresh("10.4.104.182","urn:sites:13531140:hosts:105")
	//queryallstorageunit("10.4.104.182")
	//addDatastore("10.4.104.182", "tangula_test", "urn:sites:13531140:hosts:105", "urn:sites:13531140:storageunits:85B5E76442BD4FAA9BEAA3D9913D235E")
	//getTasks("10.4.104.182")
	//delDatastore("10.4.104.182", "urn:sites:13531140:hosts:105", 22)
	queryDatastore("10.4.104.182")
	//crtVm("10.4.104.182")
	//res, _ := getVmInfo("10.4.104.182")
	//fmt.Println(util.SerializeToJson(res))
}
