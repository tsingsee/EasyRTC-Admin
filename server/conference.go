package server

import (
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr/v2"
	"tsingsee.com/adminserver/app"
	"tsingsee.com/adminserver/db"
)

// ConferenceServer 会议室服务
type ConferenceServer struct {
	*app.App
}

func NewConferenceServer(app *app.App) *ConferenceServer {
	return &ConferenceServer{
		App: app,
	}
}

// Info 获取会议室信息
func (s ConferenceServer) Info(c *gin.Context) {
	var param struct {
		ID int64 `json:"id,omitempty"`
	}
	if c.BindJSON(&param) != nil {
		return
	}
	info := app.ConferenceInfo{}
	err := s.DB().Select(app.SqlStar).From(app.ConferenceTableName).
		Where(app.WhereCommonIdAndUid, param.ID, c.GetInt64(app.UserID)).LoadOneContext(c, &info)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

// Runing 获取正在进行的会议室列表
func (s ConferenceServer) Runing(c *gin.Context) {
	items := []app.ConferenceInfo{}
	result, err := db.NewSelector(s.DB()).From(app.ConferenceTableName).Where(
		dbr.Eq(app.CommonUidCol, c.GetInt64(app.UserID)),
		dbr.Eq(app.ConferenceEtimeCol, nil),
	).LoadPage(&items)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// Dispose 解散会议室
func (s ConferenceServer) Dispose(c *gin.Context) {
	s.APIRoute(c, "/api/conference/dispose")
}

// Lock 锁定会议室
func (s ConferenceServer) Lock(c *gin.Context) {
	s.APIRoute(c, "/api/conference/lock")
}

//Unlock 解锁会议室
func (s ConferenceServer) Unlock(c *gin.Context) {
	s.APIRoute(c, "/api/conference/unlock")
}

//History 会议室历史记录
func (s ConferenceServer) History(c *gin.Context) {
	var param struct {
		RoomName string `json:"roomName,omitempty"`
		Range    struct {
			StartTime db.NullTime `json:"startTime,omitempty"`
			EndTime   db.NullTime `json:"endTime,omitempty"`
		} `json:"range,omitempty"`
		Page    uint64 `json:"page,omitempty"`
		PerPage uint64 `json:"perPage,omitempty"`
	}

	if c.BindJSON(&param) != nil {
		return
	}

	selector := db.NewSelector(s.DB())

	selector.Conditions = append(selector.Conditions, db.Condition{
		Col: app.CommonUidCol,
		Cmp: db.CmpEq,
		Val: c.GetInt64(app.UserID),
	})

	if len(param.RoomName) > 0 {
		selector.Conditions = append(selector.Conditions, db.Condition{
			Col: app.RoomNameCol,
			Cmp: db.CmpEq,
			Val: param.RoomName,
		})
	}

	if param.Range.StartTime.Valid {
		selector.Conditions = append(selector.Conditions, db.Condition{
			Col: app.CommonCtimeCol,
			Cmp: db.CmpGte,
			Val: param.Range.StartTime,
		})
	}
	if param.Range.EndTime.Valid {
		selector.Conditions = append(selector.Conditions, db.Condition{
			Col: app.CommonCtimeCol,
			Cmp: db.CmpLte,
			Val: param.Range.EndTime,
		})
	}

	selector.Orders = []db.Order{
		{Col: "id"},
	}

	confereces := []app.ConferenceInfo{}
	result, err := selector.From(app.ConferenceTableName).Paginate(param.Page, param.PerPage).LoadPage(&confereces)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

//Action 会议室事件
func (s ConferenceServer) Action(c *gin.Context) {
	req := ActionRequest{}
	if c.BindJSON(&req) != nil {
		return
	}

	if len(req.Room) == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New("room is invalid"))
		return
	}

	switch req.Action {
	case MUC_ROOM_INFO:
		logger.Info("get room.", zap.String("roomName", req.Room))
		var roomInfo app.RoomInfo
		err := s.DB().Select(app.SqlStar).From(app.RoomTableName).Where(app.WhereRoomName, req.Room).LoadOne(&roomInfo)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusOK, roomInfo)

	case MUC_ROOM_PRE_CREATE:
		logger.Info("create room.", zap.String("roomName", req.Room))
		uid, _ := s.DB().Select(app.CommonUidCol).From(app.RoomTableName).Where(app.WhereRoomName, req.Room).ReturnInt64()
		if uid == 0 {
			c.AbortWithError(http.StatusNotFound, errors.New("房间不存在"))
			return
		}
		confereceInfo := app.ConferenceInfo{
			Uid:        uid,
			RoomName:   req.Room,
			ApiEnabled: req.ApiEnabled,
			Ctime:      time.Now(),
		}
		_, err := s.DB().InsertInto(app.ConferenceTableName).
			Columns(app.CommonUidCol, app.ConferenceRoomNameCol, app.ConferenceApiEnabledCol, app.CommonCtimeCol).
			Record(&confereceInfo).ExecContext(c)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, errors.New("房间不存在"))
			return
		}
		c.JSON(http.StatusOK, confereceInfo)

	case MUC_ROOM_CREATED:
		logger.Info("created room.", zap.String("roomName", req.Room))

	case MUC_OCCUPANT_PRE_JOIN:
		participantLimits, _ := s.DB().Select(app.RoomPartLimitsCol).From(app.RoomTableName).Where(app.WhereRoomName, req.Room).ReturnInt64()
		logger.Info("pre join room.", zap.String("roomName", req.Room), zap.Int("reqLimits", req.Participants), zap.Int64("sqlLimits", participantLimits))
		if participantLimits > 0 && req.Participants >= int(participantLimits) {
			c.AbortWithError(http.StatusServiceUnavailable, errors.New("会议室人数已达上限"))
			return
		}

	case MUC_OCCUPANT_JOINED:
		logger.Info("joined room.", zap.String("roomName", req.Room))
		s.DB().Update(app.ConferenceTableName).Set(app.ConferencePartiCol, req.Participants).Where(app.WhereCommonId, req.ConferenceId).ExecContext(c)
		s.DB().Update(app.ConferenceTableName).Set(app.ConferenceMaxPartiCol, req.Participants).
			Where(app.WhereIdAndMaxParti, req.ConferenceId, req.Participants).ExecContext(c)
		// TODO: 数据库记录参会者

	case MUC_OCCUPANT_LEFT:
		logger.Info("left room.", zap.String("roomName", req.Room))
		s.DB().Update(app.ConferenceTableName).Set(app.ConferencePartiCol, req.Participants).Where(app.WhereCommonId, req.ConferenceId).ExecContext(c)
		// TODO: 数据库更新参会者

	case MUC_ROOM_DESTROYED:
		logger.Info("destory room.", zap.String("roomName", req.Room))
		s.DB().Update(app.ConferenceTableName).
			Set(app.ConferenceEtimeCol, time.Now()).
			Set(app.ConferenceIsRecordCol, false).
			Where(app.WhereCommonId, req.ConferenceId).ExecContext(c)

	case MUC_ROOM_SECRET:
		logger.Info("secret room, need password.", zap.String("roomName", req.Room))
		s.DB().Update(app.ConferenceTableName).Set(app.ConferenceLockPassCol, req.Secret).Where(app.WhereCommonId, req.ConferenceId).ExecContext(c)

	case MUC_ROOM_RECORDING_START:
		logger.Info("start recording room.", zap.String("roomName", req.Room))
		if recording := req.Recording; recording != nil {
			s.DB().Update(app.ConferenceTableName).
				Set(app.ConferenceIsRecordCol, true).
				Set(app.ConferenceStreamingCol, recording.Streaming).
				Where(app.WhereCommonId, req.ConferenceId).ExecContext(c)
		}

	case MUC_ROOM_RECORDING_STOP:
		logger.Info("stop recording room.", zap.String("roomName", req.Room))
		if recording := req.Recording; recording != nil {
			s.DB().Update(app.ConferenceTableName).
				Set(app.ConferenceIsRecordCol, false).
				Where(app.WhereCommonId, req.ConferenceId).ExecContext(c)

			uid, _ := s.DB().Select(app.CommonUidCol).From(app.RoomTableName).Where(app.WhereRoomName, req.Room).ReturnInt64()
			recordInfo := app.RecordInfo{
				Uid:          uid,
				ConferenceId: req.ConferenceId,
				RoomName:     req.Room,
				Duration:     recording.Duration,
				Size:         recording.Size,
				DownloadUrl:  recording.ObjectKey,
				StreamingUrl: recording.Streaming,
				Ctime:        time.Now(),
			}
			s.DB().InsertInto(app.RecordTableName).
				Columns(app.CommonUidCol, app.RecordConferenceIdCol, app.RecordRoomNameCol,
					app.RecordDurationCol, app.RecordSizeCol, app.RecordDownUrlCol, app.RecordStreamUrlCol, app.CommonCtimeCol).
				Record(&recordInfo).ExecContext(c)
		}
	}
}
