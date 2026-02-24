# API 文档

## 视频 TS 切片相关接口

### 保存 TS 切片
- **URL**: `/video_ts/save`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "video_id": "string",
    "key": "string",
    "iv": "string",
    "ts_data": [
      {
        "ts_sequence": "number",
        "ts_path": "string", 
        "duration": "number",
        "definition": "string"
      }
    ]
  }
  ```
- **Response**:
  ```json
  {
    "code": 1000,
    "message": "success",
    "data": {
      "video_id": "string"
    }
  }
  ```
- **错误码**:
  - `1001`: 参数绑定失败/参数验证失败
  - `1002`: 批量保存TS切片失败
  - `1003`: 保存视频加密信息失败

### 获取 TS 切片列表
- **URL**: `/video_ts/list`
- **Method**: `GET`
- **Query Parameters**:
  - `video_id`: 视频 ID（必填）
  - `definitions`: 清晰度（可选）
- **Response**:
  ```json
  {
    "code": 1000,
    "message": "success",
    "data": {
      "ts_list": [
        {
          "video_ts_id": "number",
          "video_id": "string",
          "ts_sequence": "number",
          "ts_path": "string",
          "duration": "number",
          "definition": "string",
          "create_time": "number"
        }
      ],
      "video_id": "string",
      "definitions": "string"
    }
  }
  ```
- **错误码**:
  - `1001`: 参数验证失败
  - `1002`: 查询TS切片列表失败

## 播放相关接口

### 获取播放 M3U8 文件（重定向）
- **URL**: `/play/:video_id/index.m3u8`
- **Method**: `GET`
- **Path Parameters**:
  - `video_id`: 视频 ID
- **Response**: 302 重定向到 `/play/hls/:video_id/index.m3u8`
- **错误码**:
  - `1001`: 视频ID不能为空

### 获取 HLS M3U8 文件
- **URL**: `/play/hls/:video_id/index.m3u8`
- **Method**: `GET`
- **Path Parameters**:
  - `video_id`: 视频 ID
- **Response**: M3U8 文件内容（Content-Type: application/vnd.apple.mpegurl）
  ```m3u8
  #EXTM3U
  #EXT-X-VERSION:3
  #EXT-X-TARGETDURATION:10
  #EXT-X-MEDIA-SEQUENCE:0
  #EXT-X-PLAYLIST-TYPE:VOD
  #EXT-X-KEY:METHOD=AES-128,URI="enc.key",IV=0x00000000000000000000000000000000
  #EXTINF:10.416,
  https://example.com/ts0.ts
  #EXTINF:6.833,
  https://example.com/ts1.ts
  #EXT-X-ENDLIST
  ```
- **错误码**:
  - `1001`: 视频ID不能为空
  - `1002`: 查询TS切片列表失败/获取视频加密信息失败
  - `1003`: 生成m3u8内容失败

### 获取 HLS 加密密钥
- **URL**: `/play/hls/:video_id/enc.key`
- **Method**: `GET`
- **Path Parameters**:
  - `video_id`: 视频 ID
- **Response**: 加密密钥内容（Content-Type: application/octet-stream）
- **错误码**:
  - `1001`: 视频ID不能为空
  - `1002`: 获取视频加密信息失败

## 数据实体结构

### VideoTSSaveRequest（保存TS切片请求）
```go
type VideoTSSaveRequest struct {
    VideoID string                 `json:"video_id" binding:"required"`
    TSData  []*VideoTsSaveDataItem `json:"ts_data" binding:"required"`
}
```

### VideoTsSaveDataItem（TS切片数据项）
```go
type VideoTsSaveDataItem struct {
    TSSequence int64   `json:"ts_sequence" binding:"required"`
    Key        string  `json:"key" binding:"required"`
    IV         string  `json:"iv" binding:"required"`
    TSPath     string  `json:"ts_path" binding:"required"`
    Duration   float64 `json:"duration" binding:"required"`
    Definition string  `json:"definition"`
}
```

### VideoTSEntity（TS切片实体）
```go
type VideoTSEntity struct {
    VideoTSID  int64   `gorm:"column:video_ts_id;primaryKey;autoIncrement" json:"video_ts_id"`
    VideoID    string  `gorm:"column:video_id" json:"video_id"`
    TSSequence int64   `gorm:"column:ts_sequence" json:"ts_sequence"`
    TSPath     string  `gorm:"column:ts_path" json:"ts_path"`
    Duration   float64 `gorm:"column:duration" json:"duration"`
    Definition string  `gorm:"column:definition" json:"definition"`
    CreateTime int64   `gorm:"column:create_time" json:"create_time"`
}
```

### VideoEncryptEntity（视频加密信息实体）
```go
type VideoEncryptEntity struct {
    VideoEncryptID uint64 `gorm:"column:video_encrypt_id;primaryKey;autoIncrement" json:"video_encrypt_id"`
    VideoID       string `gorm:"column:video_id;size:32;not null;index" json:"video_id"`
    Key           string `gorm:"column:key;size:64;not null" json:"key"`
    IV            string `gorm:"column:iv;size:64;not null" json:"iv"`
    CreateTime    uint64 `gorm:"column:create_time;not null;index" json:"create_time"`
}
```

## 使用示例

### 保存TS切片示例
```bash
curl -X POST http://localhost:8088/video_ts/save \
  -H "Content-Type: application/json" \
  -d '{
    "video_id": "video_123",
    "ts_data": [
      {
        "ts_sequence": 0,
        "key": "encryption_key_123",
        "iv": "initialization_vector_123",
        "ts_path": "https://example.com/ts0.ts",
        "duration": 10.416,
        "definition": "720p"
      },
      {
        "ts_sequence": 1,
        "key": "encryption_key_456",
        "iv": "initialization_vector_456",
        "ts_path": "https://example.com/ts1.ts",
        "duration": 6.833,
        "definition": "720p"
      }
    ]
  }'
```

### 获取M3U8文件示例
```bash
curl http://localhost:8088/play/hls/video_123/index.m3u8
```

### 获取加密密钥示例
```bash
curl http://localhost:8088/play/hls/video_123/enc.key
```