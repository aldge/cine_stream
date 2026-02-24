/**
 * CinePlayer JS SDK
 * 
 * 基于 DPlayer + HLS.js 的视频播放器 SDK，支持自定义加密协议解析。
 * 
 * 核心特性：
 * - 支持 AES-128-GCM 加密的 m3u8 内容解密
 * - m3u8 内容从内存读取，无需额外网络请求
 * - ts 分片由 HLS.js 原生加载，保证最佳性能
 * - 自动加载 DPlayer 和 HLS.js 依赖
 * 
 * 使用示例：
 * ```javascript
 * // 方式一：静态方法快速创建
 * const player = await CinePlayer.create('#player', 'https://api.example.com/play/video_id/index.c3u8');
 * 
 * // 方式二：实例化后初始化
 * const player = new CinePlayer();
 * await player.init('#player', 'https://api.example.com/play/video_id/index.c3u8');
 * 
 * // 控制播放
 * player.play();
 * player.pause();
 * player.seek(30);
 * 
 * // 销毁
 * player.destroy();
 * ```
 * 
 * @version 1.0.0
 * @license MIT
 */

// ==================== 常量配置 ====================

/**
 * SDK 全局配置
 * @constant {Object}
 */
const CONFIG = {
  /**
   * CDN 资源地址
   * @property {string} DPLAYER - DPlayer 播放器库
   * @property {string} HLS - HLS.js 流媒体库
   */
  CDN: {
    DPLAYER: 'https://cdn.jsdelivr.net/npm/dplayer/dist/DPlayer.min.js',
    HLS: 'https://cdn.jsdelivr.net/npm/hls.js@latest/dist/hls.min.js'
  },

  /**
   * AES-128-GCM 加密配置
   * @property {string} KEY - 16 字节密钥
   * @property {string} NONCE - 12 字节随机数
   * @property {string} ALGORITHM - 加密算法
   * @property {number} TAG_LENGTH - 认证标签长度（位）
   */
  CRYPTO: {
    KEY: '0123456789abcdef',
    NONCE: '0123456789ab',
    ALGORITHM: 'AES-GCM',
    TAG_LENGTH: 128
  },
};

// ==================== 工具函数 ====================

/**
 * Base64 字符串转 ArrayBuffer
 * @param {string} base64 - Base64 编码的字符串
 * @returns {ArrayBuffer} 解码后的二进制数据
 */
function base64ToArrayBuffer(base64) {
  const binaryString = atob(base64);
  const bytes = new Uint8Array(binaryString.length);
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }
  return bytes.buffer;
}

/**
 * 动态加载外部脚本
 * @param {string} src - 脚本 URL
 * @returns {Promise<void>} 加载完成后 resolve
 * @throws {Error} 加载失败时抛出错误
 */
function loadScript(src) {
  return new Promise((resolve, reject) => {
    const script = document.createElement('script');
    script.src = src;
    script.onload = resolve;
    script.onerror = () => reject(new Error(`Failed to load: ${src}`));
    document.head.appendChild(script);
  });
}

/**
 * 创建 HLS.js 加载统计对象
 * 用于 PlaylistLoader 回调，模拟网络请求统计信息
 * @returns {Object} HLS.js 要求的 stats 对象结构
 */
function createStats() {
  const now = performance.now();
  return {
    loading: { start: now, first: now, end: now },
    parsing: { start: now, end: now },
    buffering: { start: now, first: now, end: now },
    loaded: 0,
    total: 0
  };
}

// ==================== 加密模块 ====================

/**
 * 加密解密模块
 * 使用 Web Crypto API 实现 AES-128-GCM 解密
 */
const CryptoModule = {
  /**
   * AES-128-GCM 解密
   * @param {string} encryptedBase64 - Base64 编码的密文（包含认证标签）
   * @param {string} [key=CONFIG.CRYPTO.KEY] - 16 字节密钥
   * @returns {Promise<string>} 解密后的明文字符串
   * @throws {Error} 解密失败时抛出错误
   */
  async decrypt(encryptedBase64, key = CONFIG.CRYPTO.KEY) {
    const encoder = new TextEncoder();
    const keyBuffer = encoder.encode(key);
    const nonceBuffer = encoder.encode(CONFIG.CRYPTO.NONCE);
    const encryptedBuffer = base64ToArrayBuffer(encryptedBase64);

    // 导入密钥
    const cryptoKey = await crypto.subtle.importKey(
      'raw',
      keyBuffer,
      { name: CONFIG.CRYPTO.ALGORITHM },
      false,
      ['decrypt']
    );

    // 执行解密（自动验证认证标签）
    const plaintextBuffer = await crypto.subtle.decrypt(
      {
        name: CONFIG.CRYPTO.ALGORITHM,
        iv: nonceBuffer,
        additionalData: new Uint8Array(0),
        tagLength: CONFIG.CRYPTO.TAG_LENGTH
      },
      cryptoKey,
      encryptedBuffer
    );

    return new TextDecoder().decode(plaintextBuffer);
  }
};

// ==================== 协议解析模块 ====================

/**
 * 加密协议解析模块
 * 处理自定义加密协议的请求和响应
 */
const ProtocolModule = {
  /**
   * 检查 URL 是否为加密协议（以 .c3u8 结尾）
   * @param {string} url - 待检查的 URL
   * @returns {boolean} 是否为加密协议 URL
   */
  isEncryptedUrl(url) {
    try {
      const urlObj = new URL(url);
      return urlObj.pathname.endsWith('.c3u8');
    } catch (e) {
      return false;
    }
  },

  /**
   * 解析加密协议，获取解密后的 m3u8 内容
   * 
   * 协议响应格式：
   * ```json
   * {
   *   "code": 0,
   *   "message": "success",
   *   "data": {
   *     "info": "<Base64 编码的 AES-GCM 加密内容>"
   *   }
   * }
   * ```
   * 
   * @param {string} url - 加密协议 URL
   * @returns {Promise<string>} 解密后的 m3u8 内容
   * @throws {Error} 请求失败、协议错误或解密失败时抛出错误
   */
  async parse(url) {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`HTTP error: ${response.status}`);
    }

    const data = await response.json();
    if (data.code !== 0) {
      throw new Error(data.message || 'Protocol error');
    }

    if (!data.data?.info) {
      throw new Error('Invalid protocol format: missing data.info');
    }

    const decrypted = await CryptoModule.decrypt(data.data.info);
    if (!decrypted.includes('#EXTM3U')) {
      throw new Error('Decrypted data is not valid m3u8 format');
    }

    return decrypted;
  }
};

// ==================== Playlist Loader ====================

/**
 * 创建自定义 Playlist Loader
 * 
 * HLS.js 使用 pLoader 加载 m3u8 播放列表。
 * 此 loader 从内存读取 m3u8 内容，而非发起网络请求。
 * ts 分片仍由 HLS.js 默认 loader 加载。
 * 
 * @param {CinePlayer} player - CinePlayer 实例
 * @returns {Class} HLS.js 兼容的 Loader 类
 */
function createPlaylistLoader(player) {
  return class PlaylistLoader {
    constructor() {}

    /**
     * 加载 m3u8 内容
     * @param {Object} context - HLS.js 上下文
     * @param {Object} config - HLS.js 配置
     * @param {Object} callbacks - 回调函数集合
     */
    load(context, config, callbacks) {
      const doLoad = () => {
        callbacks.onSuccess(
          { url: context.url, data: player.m3u8Content },
          createStats(),
          context
        );
      };

      // 如果 m3u8 内容已就绪，立即返回；否则等待
      if (player.m3u8Content) {
        setTimeout(doLoad, 0);
      } else {
        player._waitForM3u8().then(() => setTimeout(doLoad, 0));
      }
    }

    abort() {}
    destroy() {}
  };
}

// ==================== 主类 ====================

/**
 * CinePlayer 视频播放器
 * 
 * 封装 DPlayer 和 HLS.js，提供简洁的 API 和加密协议支持。
 */
class CinePlayer {
  /**
   * 创建 CinePlayer 实例
   * @param {Object} [options={}] - 配置选项
   * @param {Object} [options.dpPlayerConfig={}] - DPlayer 默认配置，会与 init 时的配置合并
   */
  constructor(options = {}) {
    this._checkCompatibility();

    this.options = {
      dpPlayerConfig: {},
      ...options
    };

    /** @type {DPlayer|null} DPlayer 实例 */
    this.player = null;

    /** @type {Hls|null} HLS.js 实例 */
    this.hls = null;

    /** @type {string|null} 内存中的 m3u8 内容 */
    this.m3u8Content = null;

    /** @private Promise resolve 函数，用于等待 m3u8 就绪 */
    this._m3u8ReadyResolve = null;

    /** @private 依赖库是否已加载 */
    this._dependenciesLoaded = false;
  }

  // ==================== 私有方法 ====================

  /**
   * 检查浏览器兼容性
   * @private
   * @throws {Error} 不兼容时抛出错误
   */
  _checkCompatibility() {
    if (typeof window === 'undefined') {
      throw new Error('CinePlayer only works in browser environment');
    }
    if (!crypto?.subtle) {
      throw new Error('Browser does not support Web Crypto API');
    }
    if (typeof TextEncoder === 'undefined') {
      throw new Error('Browser does not support TextEncoder API');
    }
  }

  /**
   * 等待 m3u8 内容就绪
   * @private
   * @returns {Promise<void>}
   */
  _waitForM3u8() {
    if (this.m3u8Content) {
      return Promise.resolve();
    }
    return new Promise(resolve => {
      this._m3u8ReadyResolve = resolve;
    });
  }

  /**
   * 加载 DPlayer 和 HLS.js 依赖库
   * @private
   * @returns {Promise<void>}
   */
  async _loadDependencies() {
    if (this._dependenciesLoaded) return;

    const tasks = [];

    if (!window.DPlayer) {
      tasks.push(loadScript(CONFIG.CDN.DPLAYER));
    }
    if (!window.Hls) {
      tasks.push(loadScript(CONFIG.CDN.HLS));
    }

    if (tasks.length > 0) {
      await Promise.all(tasks);
    }

    this._dependenciesLoaded = true;
  }

  /**
   * 初始化 HLS.js 实例
   * @private
   * @param {HTMLVideoElement} video - video 元素
   * @param {string} url - 占位 URL（实际内容从内存读取）
   */
  _initHls(video, url) {
    if (!Hls.isSupported()) {
      console.error('[CinePlayer] HLS.js not supported');
      return;
    }

    // 仅自定义 pLoader 拦截 m3u8 请求
    // ts 分片使用 HLS.js 默认 loader，保证最佳性能
    this.hls = new Hls({
      pLoader: createPlaylistLoader(this),
      debug: false
    });

    this.hls.attachMedia(video);

    this.hls.on(Hls.Events.MEDIA_ATTACHED, () => {
      this.hls.loadSource(url);
    });

    this.hls.on(Hls.Events.MANIFEST_PARSED, () => {
      if (this.options.dpPlayerConfig.autoplay) {
        video.play().catch(() => {});
      }
    });

    // 错误处理与自动恢复
    this.hls.on(Hls.Events.ERROR, (_, data) => {
      if (!data.fatal) return;

      if (data.type === Hls.ErrorTypes.NETWORK_ERROR) {
        this.hls.startLoad();
      } else if (data.type === Hls.ErrorTypes.MEDIA_ERROR) {
        this.hls.recoverMediaError();
      } else {
        console.error('[CinePlayer] Fatal HLS error:', data);
      }
    });
  }

  // ==================== 公开方法 ====================

  /**
   * 手动加载 m3u8 内容到内存
   * 
   * 用于非加密协议场景，直接提供 m3u8 内容。
   * 
   * @param {string} content - m3u8 文本内容
   * 
   * @example
   * const player = new CinePlayer();
   * await player.init('#player', 'memory://playlist.m3u8');
   * player.loadM3u8Content(m3u8String);
   */
  loadM3u8Content(content) {
    this.m3u8Content = content;

    // 通知等待中的 PlaylistLoader
    if (this._m3u8ReadyResolve) {
      this._m3u8ReadyResolve();
      this._m3u8ReadyResolve = null;
    }
  }

  /**
   * 初始化播放器
   * 
   * @param {string|HTMLElement} container - 容器选择器或 DOM 元素
   * @param {string} videoUrl - 视频 URL，支持加密协议（URL 包含 /cine/）
   * @param {Object} [dpOptions={}] - DPlayer 配置，会覆盖构造函数中的默认配置
   * @returns {Promise<DPlayer>} DPlayer 实例
   * @throws {Error} 容器不存在或初始化失败时抛出错误
   * 
   * @example
   * const player = new CinePlayer();
   * await player.init('#player', 'https://api.example.com/cine/video', {
   *   autoplay: true,
   *   theme: '#00a1d6'
   * });
   */
  async init(container, videoUrl, dpOptions = {}) {
    await this._loadDependencies();

    const containerEl = typeof container === 'string'
      ? document.querySelector(container)
      : container;

    if (!containerEl) {
      throw new Error('Container not found');
    }

    // 如果不是加密协议URL，则直接使用默认的初始化方式
    if (!ProtocolModule.isEncryptedUrl(videoUrl)) {
      this.player = new DPlayer({
        ...this.options.dpPlayerConfig,
        ...dpOptions,
        container: containerEl,
        video: {
          url: videoUrl,
          type: 'hls'
        }
      });
      return this.player;
    }

    // 如果是加密协议 URL，自动解析并加载
    const placeholderUrl = 'memory://playlist.m3u8'; // 使用占位 URL，实际内容从内存读取
    this.player = new DPlayer({
      ...this.options.dpPlayerConfig,
      ...dpOptions,
      container: containerEl,
      video: {
        url: placeholderUrl,
        type: 'customHls',
        customType: {
          customHls: (video) => this._initHls(video, placeholderUrl)
        }
      }
    });
    // 使用私密协议解开
    const m3u8Content = await ProtocolModule.parse(videoUrl);
    this.loadM3u8Content(m3u8Content);
    return this.player;
  }

  /**
   * 播放视频
   */
  play() {
    this.player?.play();
  }

  /**
   * 暂停视频
   */
  pause() {
    this.player?.pause();
  }

  /**
   * 跳转到指定时间
   * @param {number} time - 目标时间（秒）
   */
  seek(time) {
    this.player?.seek(time);
  }

  /**
   * 获取当前播放时间
   * @returns {number} 当前时间（秒）
   */
  get currentTime() {
    return this.player?.video?.currentTime || 0;
  }

  /**
   * 获取视频总时长
   * @returns {number} 总时长（秒）
   */
  get duration() {
    return this.player?.video?.duration || 0;
  }

  /**
   * 获取 DPlayer 实例
   * 可用于访问 DPlayer 的完整 API
   * @returns {DPlayer|null}
   */
  get dp() {
    return this.player;
  }

  /**
   * 销毁播放器，释放资源
   * 销毁后实例不可再使用
   */
  destroy() {
    if (this.hls) {
      this.hls.destroy();
      this.hls = null;
    }

    if (this.player) {
      this.player.destroy();
      this.player = null;
    }

    this.m3u8Content = null;
    this._m3u8ReadyResolve = null;
  }

  /**
   * 静态工厂方法：快速创建并初始化播放器
   * 
   * @param {string|HTMLElement} container - 容器选择器或 DOM 元素
   * @param {string} url - 视频 URL
   * @param {Object} [options={}] - 配置选项（同时作为构造函数和 init 的参数）
   * @returns {Promise<CinePlayer>} CinePlayer 实例
   * 
   * @example
   * const player = await CinePlayer.create('#player', 'https://api.example.com/cine/video', {
   *   dpPlayerConfig: { autoplay: true }
   * });
   */
  static async create(container, url, options = {}) {
    const instance = new CinePlayer(options);
    await instance.init(container, url, options);
    return instance;
  }
}

// ==================== 导出 ====================

// 浏览器全局变量
if (typeof window !== 'undefined') {
  window.CinePlayer = CinePlayer;
}

// CommonJS 模块导出
if (typeof module !== 'undefined' && module.exports) {
  module.exports = CinePlayer;
}
