### 获取验证码id
POST http://localhost:8004/admin/captcha-id
Accept: */*
Cache-Control: no-cache
Content-Type: application/json

### 获取图片验证码
GET http://localhost:8004/admin/captcha/2uVP5dd1mOIw0ShdNNc7.png
Accept: */*
Cache-Control: no-cache

### 获取语音验证码
GET http://localhost:8004/admin/captcha/DEpEyBEkn1tRq5rCB5xy.wav
Accept: */*
Cache-Control: no-cache

### 注册用户
POST http://localhost:8004/admin/passport/signup
Accept: */*
Cache-Control: no-cache
Content-Type: application/json

{
  "name": "14131913",
  "password": "123456",
  "captcha_id": "XgxjDNvGFIuBjUYfWggb",
  "captcha_code": "676832"
}

### 用户登录
POST http://localhost:8004/admin/passport/login
Accept: */*
Cache-Control: no-cache
Content-Type: application/json

{
  "name": "14131913",
  "password": "123456",
  "captcha_id": "2uVP5dd1mOIw0ShdNNc7",
  "captcha_code": "472991"
}

### 用户注销
POST http://localhost:8004/admin/passport/logout
Accept: */*
Cache-Control: no-cache
Content-Type: application/json

### 获取个人账户信息
POST http://localhost:8004/admin/passport/info
Accept: */*
Cache-Control: no-cache
Content-Type: application/json

### 修改个人账户信息
POST http://localhost:8004/admin/passport/modify
Accept: */*
Cache-Control: no-cache
Content-Type: application/json

{
  "displayName": "显示1名",
  "password": "123456",
  "newpass": "456789",
  "email": "459685578@qq.com",
  "company": "安徽旭帆",
  "phone": "1825543957"
}

### 创建房间
POST http://localhost:8004/admin/room/create
Accept: */*
Cache-Control: no-cache
Content-Type: application/json

{
  "roomName": "测试房间3",
  "participantLimits": 10,
  "allowAnonymous": false,
  "roomConfig": {
    "resolution": 320,
    "subject": "房间主题",
    "lockPassword": "8384",
    "requireDisplayName": true,
    "startWithAudioMuted": false,
    "startWithVideoMuted": false,
    "fileRecordingsEnabled": false,
    "liveStreamingEnabled": false,
    "bandwidth": 1000
  }
}

### 获取房间
POST http://localhost:8004/admin/room/info
Accept: */*
Cache-Control: no-cache
Content-Type: application/json
Cookie: rtcadmin=test

{
  "id": 2
}

### 修改房间
POST http://localhost:8004/admin/room/modify
Accept: */*
Cache-Control: no-cache
Content-Type: application/json
Cookie: rtcadmin=test

{
  "id": 1,
  "roomName": "测试房间",
  "participantLimits": 13,
  "allowAnonymous": false,
  "roomConfig": {
    "resolution": 320,
    "subject": "房间配置主题",
    "lockPassword": "8384",
    "requireDisplayName": true,
    "startWithAudioMuted": false,
    "startWithVideoMuted": false,
    "fileRecordingsEnabled": false,
    "liveStreamingEnabled": false,
    "bandwidth": 1000
  }
}

### 获取房间列表
POST http://localhost:8004/admin/room/list
Accept: */*
Cache-Control: no-cache
Content-Type: application/json
Cookie: rtcadmin=test

{
  "page": 0,
  "perPage": 10
}

### 删除房间
POST http://localhost:8004/admin/room/delete
Accept: */*
Cache-Control: no-cache
Content-Type: application/json
Cookie: rtcadmin=test

{
  "id": 1
}

### 获取 token
POST http://localhost:8004/admin/room/token

### 获取视频回看
POST http://localhost:8004/admin/record/info
Accept: */*
Cache-Control: no-cache
Content-Type: application/json
Cookie: rtcadmin=test

{
  "id": 1
}

### 获取视频回看列表
POST http://localhost:8004/admin/record/list
Accept: */*
Cache-Control: no-cache
Content-Type: application/json
Cookie: rtcadmin=test

{
  "page": 0,
  "perPage": 10
}

### 删除视频回看
POST http://localhost:8004/admin/record/delete
Accept: */*
Cache-Control: no-cache
Content-Type: application/json
Cookie: rtcadmin=test

{
  "id": 1
}



