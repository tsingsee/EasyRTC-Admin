package app

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"jhmeeting.com/adminserver/db"
)

const (
	UserID = "uid"
)

const (
	SqlStar             = "*"
	CommonCtimeCol      = "ctime"
	CommonIdCol         = "id"
	CommonUidCol        = UserID
	WhereCommonId       = "id=?"
	WhereCommonIdAndUid = "id=? and uid=?"
)

//*****************************************用户数据*********************************************************/
// 用户
type User struct {
	Id          int64     `json:"id,omitempty"`       // id
	Name        string    `json:"name,omitempty"`     // 登录名
	Password    string    `json:"password,omitempty"` // 密码
	DisplayName string    `json:"displayName"`        // 姓名
	Email       string    `json:"email"`              // 邮箱
	Phone       string    `json:"phone"`              // 手机号码
	Company     string    `json:"company"`            // 公司名称
	Ctime       time.Time `json:"ctime,omitempty"`    // 创建时间
}

// 用户表对应的表名称和字段名称
const (
	UserTableName   = "users"
	UserNameCol     = "name"
	UserPasswordCol = "password"
	UserDisNameCol  = "display_name"
	UserEmailCol    = "email"
	UserPhoneCol    = "phone"
	UserCompanyCol  = "company"
	WhereUserName   = "name=?"
)

//*****************************************用户创建会议室*********************************************************/
// 房间信息
type RoomInfo struct {
	Id                int64      `json:"id,omitempty"`
	Uid               int64      `json:"uid,omitempty" sql:"index:ri_uid"`      // 房间uid
	RoomName          string     `json:"roomName" sql:"index:ri_room_name"`     // 房间名称
	ParticipantLimits int        `json:"participantLimits"`                     // 房间最高参会人数
	AllowAnonymous    bool       `json:"allowAnonymous"`                        // 是否允许匿名用户创建会议
	Config            RoomConfig `json:"roomConfig,omitempty"`                  // 房间配置
	Ctime             time.Time  `json:"ctime,omitempty"  sql:"index:ri_ctime"` // 创建时间
}

// 房间表对应的表名称和字段名称
const (
	RoomTableName         = "room"
	RoomNameCol           = "room_name"
	RoomPartLimitsCol     = "participant_limits"
	RoomAllowAnonymousCol = "allow_anonymous"
	RoomConfigCol         = "config"
	WhereRoomName         = "room_name=?"
)

// 房间配置
type RoomConfig struct {
	Resolution            int    `json:"resolution,omitempty"`  // 分辨率，360，480，720，1080
	Subject               string `json:"subject"`               // 主题
	LockPassword          string `json:"lockPassword"`          // 进入密码
	RequireDisplayName    bool   `json:"requireDisplayName"`    // 是否提示参会者输入名字
	StartWithAudioMuted   bool   `json:"startWithAudioMuted"`   // 是否加入会议室时不开启音频
	StartWithVideoMuted   bool   `json:"startWithVideoMuted"`   // 是否加入会议室时不开启视频
	FileRecordingsEnabled *bool  `json:"fileRecordingsEnabled"` // 是否允许服务器录制
	LiveStreamingEnabled  *bool  `json:"liveStreamingEnabled"`  // 是否允许直播
	Bandwidth             int    `json:"bandwidth"`             // 设置参会者最大上行比特率，默认各分辨率对应的比特率，360：800，480：1000，720：1500，1080：3000
}

func (config RoomConfig) Value() (driver.Value, error) {
	data, _ := json.Marshal(config)

	return string(data), nil
}

func (config *RoomConfig) Scan(src interface{}) error {
	var source []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return errors.New("Incompatible type for RoomConfig")
	}

	return json.Unmarshal(source, config)
}

//*****************************************正在开会会议室*********************************************************/
// 会议室信息，会议室表示正在开会的房间
type ConferenceInfo struct {
	Id              int64       `json:"id,omitempty"`
	Uid             int64       `json:"uid,omitempty" sql:"index:ci_uid"`            // 会议uid
	RoomName        string      `json:"roomName,omitempty" sql:"index:ci_room_name"` // 房间名称
	Participants    int         `json:"participants,omitempty"`                      // 当前人数
	MaxParticipants int         `json:"maxParticipants,omitempty"`                   // 最高人数
	IsRecording     bool        `json:"isRecording,omitempty"`                       // 是否正在录制，直播也是录制
	Streaming       string      `json:"streaming,omitempty"`                         // 直播地址，录制则需清空
	ApiEnabled      bool        `json:"apiEnabled,omitempty"`                        // 是否是使用API接入的会议室
	LockPassword    string      `json:"lockPassword,omitempty"`                      // 进入密码
	Locked          bool        `json:"locked,omitempty"`                            // 是否锁定
	Ctime           time.Time   `json:"ctime,omitempty" sql:"index:ci_ctime"`        // 开始时间
	Etime           db.NullTime `json:"etime,omitempty" sql:"index:ci_etime"`        // 结束时间
}

// 房间表对应的表名称和字段名称
const (
	ConferenceTableName     = "conference"
	ConferenceEtimeCol      = "etime"
	ConferenceRoomNameCol   = "room_name"
	ConferenceApiEnabledCol = "api_enabled"
	ConferencePartiCol      = "participants"
	ConferenceMaxPartiCol   = "max_participants"
	ConferenceIsRecordCol   = "is_recording"
	ConferenceStreamingCol  = "streaming"
	ConferenceLockPassCol   = "lock_password"

	WhereIdAndMaxParti = "id=? and max_participants<?"
)

//*****************************************会议回看定义*********************************************************/
// 会议回看信息
type RecordInfo struct {
	Id           int64     `json:"id,omitempty"`
	ConferenceId int64     `json:"conferenceId,omitempty"` // 会议室id
	Uid          int64     `json:"uid,omitempty"`          // 会议室用户id
	RoomName     string    `json:"roomName,omitempty"`     // 会议室名称
	Duration     int64     `json:"duration,omitempty"`     // 录制时长
	Size         int64     `json:"size,omitempty"`         // 文件大小
	DownloadUrl  string    `json:"downloadUrl,omitempty"`  // 录像 url 地址
	StreamingUrl string    `json:"streamingUrl,omitempty"` // 推流 url 地址
	Ctime        time.Time `json:"ctime,omitempty"`        // 开始时间
}

// 会议回看表对应的字符串
const (
	RecordTableName       = "record"
	RecordConferenceIdCol = "conference_id"
	RecordRoomNameCol     = "room_name"
	RecordDurationCol     = "duration"
	RecordSizeCol         = "size"
	RecordDownUrlCol      = "download_url"
	RecordStreamUrlCol    = "streaming_url"

	WhereRecordConfIDAndStream = "conference_id=? and streaming_url=? and duration=0"
)
