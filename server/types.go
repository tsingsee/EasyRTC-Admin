package server

import (
	"tsingsee.com/adminserver/util"
	"time"
)

var logger = util.GetLogger()

const (
	MUC_ROOM_PRE_CREATE      = "muc-room-pre-create"      // 创建会议前事件
	MUC_ROOM_CREATED         = "muc-room-created"         // 创建会议后事件
	MUC_OCCUPANT_PRE_JOIN    = "muc-occupant-pre-join"    // 加入会议前事件
	MUC_OCCUPANT_JOINED      = "muc-occupant-joined"      // 加入会议后事件
	MUC_OCCUPANT_LEFT        = "muc-occupant-left"        // 离开会议事件
	MUC_ROOM_DESTROYED       = "muc-room-destroyed"       // 结束会议事件
	MUC_ROOM_SECRET          = "muc-room-secret"          // 设置会议室密码事件
	MUC_ROOM_INFO            = "muc-room-info"            // 获取房间信息
	MUC_ROOM_RECORDING_START = "muc-room-recording-start" // 开始录制事件
	MUC_ROOM_RECORDING_STOP  = "muc-room-recording-stop"  // 结束录制事件
)

type ActionRequest struct {
	Action       string         `json:"action,omitempty"`       // 事件名
	ConferenceId int64          `json:"conferenceId,omitempty"` // 会议室ID
	Room         string         `json:"room,omitempty"`         // 房间名
	Nick         string         `json:"nick,omitempty"`         // 参会者昵称
	Jid          string         `json:"jid,omitempty"`          // 参会者ID
	Secret       string         `json:"secret,omitempty"`       // 会议室密码
	Participants int            `json:"participants,omitempty"` // 参会人数
	ApiEnabled   bool           `json:"apiEnabled,omitempty"`   // 是否使用SDK接入
	Recording    *RecordingFile `json:"recording,omitempty"`    // 录制文件
}

type RecordingFile struct {
	ObjectKey string    `json:"objectKey,omitempty"` // 服务器录制文件地址，具体格式由存储方式确定
	Room      string    `json:"room,omitempty"`      // 房间名
	Size      int64     `json:"size,omitempty"`      // 录制文件大小（bytes）
	Duration  int64     `json:"duration,omitempty"`  // 房间名
	Streaming string    `json:"streaming,omitempty"` // 推流地址
	Ctime     time.Time `json:"ctime,omitempty"`     // 创建时间
}

type RoomTokenRequest struct {
	RoomName  string   `json:"roomName,omitempty" binding:"required"`
	ExpiresAt int64    `json:"expiresAt,omitempty"`
	Context   *Context `json:"context,omitempty"`
	Anonymous bool     `json:"anonymous,omitempty"`
}

type Context struct {
	User   *ContextUserInfo `json:"user,omitempty"`
	Callee *ContextUserInfo `json:"callee,omitempty"`
}

type ContextUserInfo struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	AvatarUrl string `json:"avatarUrl,omitempty"`
}
