package response

import (
	"gopher/entity/activity/activitymodel"
	"gopher/entity/activity/activityrepo"
	"gopher/entity/service"
)

// RecordCreateInstant make it simpler for calling the record
func (r *Response) RecordCreateInstant(event string, newData interface{}) {
	r.Record(event, nil, newData)
}

// RecordInstant is used for saving activity
// TODO: deprecated
func (r *Response) RecordInstant(event string, data ...interface{}) {
	activityServ := service.ProvideActivityService(activityrepo.ProvideActivityRepo(r.Engine))
	activityServ.Record(r.Context, event, data...)
}

// Record will send the activity for read/update/delete to the AcitivityCh
func (r *Response) Record(event string, data ...interface{}) {
	r.initiateRecordCh(event, data...)
}

// RecordCreate will send the activity for creation to the AcitivityCh
func (r *Response) RecordCreate(event string, newData interface{}) {
	r.initiateRecordCh(event, nil, newData)
}

func (r *Response) initiateRecordCh(event string, data ...interface{}) {
	activityServ := service.ProvideActivityService(activityrepo.ProvideActivityRepo(r.Engine))

	var userID uint
	var username string

	recordType := activityServ.FindRecordType(data...)
	before, after := activityServ.FillBeforeAfter(recordType, data...)

	if len(data) > 0 && !r.Engine.Environments.Activity.Write {
		return
	}

	if len(data) == 0 && !r.Engine.Environments.Activity.Read {
		return
	}

	if activityServ.IsRecordSetInEnvironment(recordType) {
		return
	}

	if userIDtmp, ok := r.Context.Get("USER_ID"); ok {
		userID = userIDtmp.(uint)
	}
	if usernameTmp, ok := r.Context.Get("USERNAME"); ok {
		username = usernameTmp.(string)
	}

	activity := activitymodel.Activity{
		Event:    event,
		UserID:   userID,
		Username: username,
		IP:       r.Context.Request.Header.Get("X-User-IP"), //r.Context.ClientIP(),
		URI:      r.Context.Request.RequestURI,
		Before:   string(before),
		After:    string(after),
	}

	r.Engine.ActivityCh <- activity

	_ = activity
	// activityServ.RecordCh(ac

}
