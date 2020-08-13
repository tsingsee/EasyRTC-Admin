package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tsingsee.com/adminserver/app"
	"tsingsee.com/adminserver/db"
)

type RecordServer struct {
	*app.App
}

func NewRecordServer(app *app.App) *RecordServer {
	return &RecordServer{
		App: app,
	}
}

func (s RecordServer) Info(c *gin.Context) {
	var param struct {
		ID int64
	}
	if c.BindJSON(&param) != nil {
		return
	}

	uid := c.GetInt64(app.UserID)

	record := app.RecordInfo{}
	err := s.DB().Select(app.SqlStar).From(app.RecordTableName).
		Where(app.WhereCommonIdAndUid, param.ID, uid).LoadOneContext(c, &record)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, record)
}

// 删除
func (s RecordServer) Delete(c *gin.Context) {
	var param struct {
		ID int64
	}
	if c.BindJSON(&param) != nil {
		return
	}
	uid := c.GetInt64(app.UserID)

	_, err := s.DB().
		DeleteFrom(app.RecordTableName).Where(app.WhereCommonIdAndUid, param.ID, uid).ExecContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

// 列表查看
func (s RecordServer) List(c *gin.Context) {
	var param struct {
		RoomName string `json:"roomName,omitempty"`
		Page     uint64 `json:"page,omitempty"` // start from 0
		PerPage  uint64 `json:"perPage,omitempty"`
	}
	//var param db.Pagination
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

	records := []*app.RecordInfo{}

	result, _ := selector.From(app.RecordTableName).
		Paginate(param.Page, param.PerPage).
		OrderDesc(app.CommonIdCol).
		LoadPage(&records)

	for _, record := range records {
		record.DownloadUrl = s.Config().RecordingURL + record.DownloadUrl
	}

	c.JSON(http.StatusOK, result)
}
