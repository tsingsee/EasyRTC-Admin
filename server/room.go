package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocraft/dbr/v2"
	"tsingsee.com/adminserver/app"
	"tsingsee.com/adminserver/db"
)

type RoomServer struct {
	*app.App
}

func NewRoomServer(app *app.App) *RoomServer {
	return &RoomServer{
		App: app,
	}
}

func (s RoomServer) Info(c *gin.Context) {
	var param struct {
		ID int64
	}
	if c.BindJSON(&param) != nil {
		return
	}
	uid := c.GetInt64(app.UserID)
	room := app.RoomInfo{}
	err := s.DB().Select(app.SqlStar).From(app.RoomTableName).
		Where(app.WhereCommonIdAndUid, param.ID, uid).LoadOneContext(c, &room)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, room)
}

func (s RoomServer) Create(c *gin.Context) {
	roomInfo := app.RoomInfo{}
	if c.BindJSON(&roomInfo) != nil {
		return
	}

	roomInfo.Uid = c.GetInt64(app.UserID)
	roomInfo.Ctime = time.Now()

	room := app.RoomInfo{}
	err := s.DB().Select(app.SqlStar).From(app.RoomTableName).
		Where(app.WhereRoomName, roomInfo.RoomName).LoadOneContext(c, &room)
	if room.RoomName == roomInfo.RoomName {
		c.AbortWithError(http.StatusBadRequest, errors.New("会议名已存在（会议名全部唯一）！"))
		return
	}

	_, err = s.DB().InsertInto(app.RoomTableName).
		Columns(app.CommonUidCol, app.RoomPartLimitsCol, app.RoomNameCol, app.RoomAllowAnonymousCol, app.RoomConfigCol, app.CommonCtimeCol).
		Record(&roomInfo).ExecContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": roomInfo.Id,
	})
}

func (s RoomServer) Delete(c *gin.Context) {
	var param struct {
		ID int64
	}
	if c.BindJSON(&param) != nil {
		return
	}
	uid := c.GetInt64(app.UserID)
	_, err := s.DB().DeleteFrom(app.RoomTableName).Where(app.WhereCommonIdAndUid, param.ID, uid).ExecContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (s RoomServer) Modify(c *gin.Context) {
	roomInfo := app.RoomInfo{}
	if c.BindJSON(&roomInfo) != nil {
		return
	}
	uid := c.GetInt64(app.UserID)
	_, err := s.DB().Update(app.RoomTableName).
		Set(app.RoomPartLimitsCol, roomInfo.ParticipantLimits).
		Set(app.RoomAllowAnonymousCol, roomInfo.AllowAnonymous).
		Set(app.RoomConfigCol, roomInfo.Config).
		Where(app.WhereCommonIdAndUid, roomInfo.Id, uid).
		ExecContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (s RoomServer) List(c *gin.Context) {
	var param db.Pagination
	if c.BindJSON(&param) != nil {
		return
	}

	uid := c.GetInt64(app.UserID)
	rooms := []app.RoomInfo{}

	result, _ := db.NewSelector(s.DB()).From(app.RoomTableName).
		Where(dbr.Eq(app.CommonUidCol, uid)).
		Paginate(param.Page, param.PerPage).
		OrderDesc(app.CommonIdCol).
		LoadPage(&rooms)
	c.JSON(http.StatusOK, result)
}

func (s RoomServer) Token(c *gin.Context) {
	// JSON Body: see RoomTokenRequest
	s.APIRoute(c, "/api/conference/token")
}
