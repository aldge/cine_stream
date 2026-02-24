## CinePlayer JS SDK

基于 DPlayer + HLS.js 的视频播放器 SDK，支持 AES-128-GCM 加密的 m3u8 内容解密。

### 核心特性

- ✅ 基于 DPlayer + HLS.js 实现
- ✅ 支持 AES-128-GCM 加密的 m3u8 内容解密
- ✅ m3u8 内容从内存读取，无需额外网络请求
- ✅ ts 分片由 HLS.js 原生加载，保证最佳性能
- ✅ 自动加载 DPlayer 和 HLS.js 依赖
- ✅ 支持以 `.c3u8` 结尾的加密协议 URL
- ✅ 提供简洁 API 和多种初始化方式

### 安装和使用

#### 1. 安装依赖

```bash
npm install
```

#### 2. 构建 SDK

```bash
npm run build     # 生产构建 -> dist/cine-player.min.js (混淆压缩)
npm run build:dev # 开发构建 -> dist/cine-player.js
```

#### 3. 启动本地服务器

```bash
npm start
```

服务器将在 http://localhost:8080 启动

测试地址：http://localhost:8080/play.html?url=https://api.example.com/play/video_id/index.c3u8

### SDK API

#### 静态方法（推荐）

```javascript
// 快速创建并初始化播放器
const player = await CinePlayer.create('#player', 'https://api.example.com/play/video_id/index.c3u8', {
  autoplay: true,
  theme: '#00a1d6'
});
```

#### 实例化方式

```javascript
// 创建实例
const player = new CinePlayer({
  dpPlayerConfig: {
    autoplay: true,
    controls: true
  }
});

// 初始化播放器
await player.init('#player', 'https://api.example.com/play/video_id/index.c3u8', {
  theme: '#ff6b35'
});

// 控制播放
player.play();
player.pause();
player.seek(30);

// 获取播放状态
console.log(player.currentTime);
console.log(player.duration);

// 销毁播放器
player.destroy();
```

### 加密协议格式

SDK 支持以下格式的加密协议响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "info": "<Base64 编码的 AES-GCM 加密内容>"
  }
}
```

解密后的内容必须是有效的 m3u8 格式（包含 `#EXTM3U`）。

### 使用示例

#### 1. 在 HTML 中使用

```html
<!DOCTYPE html>
<html>
<head>
    <title>CinePlayer Demo</title>
</head>
<body>
    <div id="player"></div>
    <script src="dist/cine-player.min.js"></script>
    <script>
        CinePlayer.create('#player', 'https://api.example.com/video.c3u8', {
            width: '100%',
            height: '600px',
            autoplay: true
        });
    </script>
</body>
</html>
```

#### 2. 通过 URL 参数播放

访问：`http://127.0.0.1:8080/play.html?url=https://api.example.com/play/video_id/index.c3u8`

#### 3. 手动输入 URL 播放

打开 `http://127.0.0.1:8080/play.html`，在输入框中粘贴以 `.c3u8` 结尾的加密协议 URL，点击播放。

### 项目结构

```
cine_player/
├── src/
│   └── cine-player.js          # SDK 核心代码
├── public/
│   ├── play.html             # 播放器页面（可独立嵌入）
│   └── dist/
│       ├── cine-player.js      # 开发版构建产物
│       └── cine-player.min.js  # 生产版构建产物（混淆压缩）
├── package.json                 # 项目配置
├── webpack.config.js           # 构建配置
└── README.md                   # 说明文档
```

### 协议解析流程

1. 检测 URL 路径是否以 `.c3u8` 结尾
2. 发送请求获取加密协议数据
3. Base64 解码 `data.info` 字段
4. AES-128-GCM 解密获取真实 m3u8 内容
5. 将 m3u8 内容加载到内存，HLS.js 从内存读取播放列表
6. ts 分片由 HLS.js 原生加载和播放

### 注意事项

- 浏览器需支持 Web Crypto API
- 确保 DPlayer 和 HLS.js 能正常从 CDN 加载
- 加密协议必须返回正确的 code=0 表示成功
- 支持跨域请求，需确保服务端配置 CORS
- 解密后的内容必须是有效的 m3u8 格式
- player.html 设计为可独立嵌入其他网站

### 浏览器支持

- Chrome >= 60
- Firefox >= 55
- Safari >= 12
- Edge >= 79