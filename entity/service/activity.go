package service

import (
	"encoding/json"
	"gopher/entity/activity/activitymodel"
	"gopher/entity/activity/activityrepo"
	"gopher/internal/core"
	"gopher/internal/param"
	"time"

	"github.com/gin-gonic/gin"
)

// RecordType is and int used as an enum
type RecordType int

const (
	read RecordType = iota
	writeBefore
	writeAfter
	writeBoth
)

// ActivityService for injecting auth
type ActivityService struct {
	Repo   activityrepo.ActivityRepo
	Engine *core.Engine
}

// ProvideActivityService for activity is used in wire
func ProvideActivityService(p activityrepo.ActivityRepo) ActivityService {
	return ActivityService{Repo: p, Engine: p.Engine}
}

// Save activity
func (s *ActivityService) Save(activity activitymodel.Activity, params param.Param) (createdActivity activitymodel.Activity, err error) {
	createdActivity, err = s.Repo.Create(activity, params)
	return
}

// ActivityWatcher is used for watching activity channel
func (s *ActivityService) ActivityWatcher() {
	var arr []activitymodel.Activity
	var counter uint64
	counter = 0
	var activity activitymodel.Activity

	tickTimer := time.Tick(time.Duration(s.Engine.Environments.Activity.SaveAfter) * time.Second)

	for {
		select {
		case activity = <-s.Engine.ActivityCh:
			counter++
			arr = append(arr, activity)
			if counter > s.Engine.Environments.Activity.BatchSize {
				s.Repo.CreateBatch(arr, param.Param{})
				counter = 0
				arr = []activitymodel.Activity{}
			}
		case <-tickTimer:
			if len(arr) > 0 {
				s.Repo.CreateBatch(arr, param.Param{})
				counter = 0
				arr = []activitymodel.Activity{}
			}
		}
	}
}

// Record will save the activity
// TODO: Record is deprecated we should go with channels
func (s *ActivityService) Record(c *gin.Context, event string, data ...interface{}) {
	var userID uint
	var username string

	recordType := s.FindRecordType(data...)
	before, after := s.FillBeforeAfter(recordType, data...)

	if len(data) > 0 {
		return
	}

	if len(data) == 0 {
		return
	}

	if userIDtmp, ok := c.Get("USER_ID"); ok {
		userID = userIDtmp.(uint)
	}
	if usernameTmp, ok := c.Get("USERNAME"); ok {
		username = usernameTmp.(string)
	}

	activity := activitymodel.Activity{
		Event:    event,
		UserID:   userID,
		Username: username,
		IP:       c.Request.Header.Get("X-User-IP"), //r.Context.ClientIP(),
		URI:      c.Request.RequestURI,
		Before:   string(before),
		After:    string(after),
	}

	_, err := s.Repo.Create(activity, param.Param{})
	s.Engine.ServerLog.CheckError(err, "E1000080", "Failed in saving activity", activity)
}

// FillBeforeAfter check if there is a need for entering before data or not
func (s *ActivityService) FillBeforeAfter(recordType RecordType, data ...interface{}) (before, after []byte) {
	var err error
	if recordType == writeBefore || recordType == writeBoth {
		before, err = json.Marshal(data[0])
		s.Engine.ServerLog.CheckError(err, "E1000081", "error in encoding data to before-json")
	}
	if recordType == writeAfter || recordType == writeBoth {
		after, err = json.Marshal(data[1])
		s.Engine.ServerLog.CheckError(err, "E1000082", "error in encoding data to after-json")
	}

	return
}

// FindRecordType is helper function for finding the best way for recording data
func (s *ActivityService) FindRecordType(data ...interface{}) RecordType {
	switch len(data) {
	case 0:
		return read
	case 2:
		return writeBoth
	default:
		if data[0] == nil {
			return writeAfter
		}
	}

	return writeBefore
}

// IsRecordSetInEnvironment check if in the env file record activated or not
func (s *ActivityService) IsRecordSetInEnvironment(recordType RecordType) bool {
	switch recordType {
	case read:
		if !s.Engine.Environments.Activity.Read {
			return true
		}
	default:
		if !s.Engine.Environments.Activity.Write {
			return true
		}
	}
	return false
}

// List of activities, it support pagination and search and return back count
func (s *ActivityService) List(params param.Param) (data map[string]interface{}, err error) {

	data = make(map[string]interface{})

	data["list"], err = s.Repo.List(params)
	s.Engine.ServerLog.CheckError(err, "E1000083", "activities list")
	if err != nil {
		return
	}

	data["count"], err = s.Repo.Count(params)
	s.Engine.ServerLog.CheckError(err, "E1000084", "activities count")

	return
}
