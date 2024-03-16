package logmgm

import (
	"github.com/Zhang-jie-jun/tangula/contants"
	"github.com/Zhang-jie-jun/tangula/internal/dao/log"
	"github.com/Zhang-jie-jun/tangula/pkg/errors"
	"github.com/Zhang-jie-jun/tangula/pkg/msg"
	"github.com/Zhang-jie-jun/tangula/pkg/util"
	"github.com/Zhang-jie-jun/tangula/routers/view"
	"github.com/sirupsen/logrus"
	"time"
)

func GetLogRecord(queryParam *view.LogRecordQueryParam) (totalNum int64, result []map[string]interface{}, err error) {
	totalNum, records, err := log.LogMgm.GetHostListByUser(queryParam.Index, queryParam.Count, queryParam.User)
	if err != nil {
		err = errors.New(msg.ERROR_GET_LOG_RECORD_INFO, msg.GetMsg(msg.ERROR_GET_LOG_RECORD_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	for _, iter := range records {
		result = append(result, iter.TransformMap())
	}
	return
}

func CreateLogRecord(operation, object, detail, user string, status contants.LogStatus) (logRecord log.LogRecord, err error) {
	if operation == "" {
		operation = msg.GetOperation(msg.UNDEFINE_OPERAT)
	}
	if object == "" {
		object = "unknown"
	}
	var record log.LogRecord
	record.Operation = operation
	record.Object = object
	record.Detail = detail
	record.User = user
	record.Status = status
	record.CreateTime = util.Time{Time: time.Now()}
	logRecord, err = log.LogMgm.CreateRecord(record)
	if err != nil {
		err = errors.New(msg.ERROR_CREATE_LOG_RECORD_INFO, msg.GetMsg(msg.ERROR_CREATE_LOG_RECORD_INFO, err.Error()))
		logrus.Error(err)
		return
	}
	return
}
