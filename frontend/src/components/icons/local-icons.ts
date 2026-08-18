// 内置本地 SVG 图标库（离线可用，不依赖外部字体）
// 每个图标使用 24x24 viewBox 的 path 数据，统一 currentColor 着色
export interface LocalIconDef {
  name: string;
  category: string;
  path: string;
}

// 通用 24x24 几何/业务图标 path 集
export const localIcons: LocalIconDef[] = [
  // 基础
  { name: 'home', category: '基础', path: 'M3 11l9-8 9 8v9a1 1 0 0 1-1 1h-5v-6H9v6H4a1 1 0 0 1-1-1z' },
  { name: 'user', category: '基础', path: 'M12 12a5 5 0 1 0 0-10 5 5 0 0 0 0 10zm0 2c-4 0-8 2-8 5v1h16v-1c0-3-4-5-8-5z' },
  { name: 'users', category: '基础', path: 'M9 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8zm7 0a4 4 0 1 0 0-8 4 4 0 0 0 0 8zm-7 2c-4 0-7 2-7 5v1h10v-1c0-1.5 1-3 3-4-1.5-.6-2-1-3-1H9zm14 4c0-3-3-5-7-5-.6 0-1.2.05-1.7.15C15.8 13 17 14.5 17 16v1h6z' },
  { name: 'settings', category: '基础', path: 'M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8zm9 4l-2 1.5.6 2.5-2.3 1-1.7 1.9-2.6-.4-1.7 1.7H10.7l-1.7-1.7-2.6.4L4.7 17l-.6-2.5L2 13l1.4-1.5L2.8 9l2.3-1 1.7-1.9 2.6.4 1.7-1.7h2.6l1.7 1.7 2.6-.4 1.7 1.9 2.3 1-.6 2.5z' },
  { name: 'search', category: '基础', path: 'M10 4a6 6 0 1 0 3.7 10.7l4.3 4.3 1.4-1.4-4.3-4.3A6 6 0 0 0 10 4zm0 2a4 4 0 1 1 0 8 4 4 0 0 1 0-8z' },
  { name: 'bell', category: '基础', path: 'M12 2a6 6 0 0 0-6 6v4l-2 3v1h16v-1l-2-3V8a6 6 0 0 0-6-6zm0 20a3 3 0 0 0 3-3H9a3 3 0 0 0 3 3z' },
  { name: 'star', category: '基础', path: 'M12 2l3 6.3 6.9 1-5 4.9 1.2 6.8L12 17.8 5.9 21l1.2-6.8-5-4.9 6.9-1z' },
  { name: 'heart', category: '基础', path: 'M12 21s-7-4.5-9.5-9C1 9 2.5 5 6 5c2 0 3.2 1.2 4 2.5C10.8 6.2 12 5 14 5c3.5 0 5 4 3.5 7-2.5 4.5-9.5 9-9.5 9z' },
  { name: 'bookmark', category: '基础', path: 'M6 3h12a1 1 0 0 1 1 1v17l-7-4-7 4V4a1 1 0 0 1 1-1z' },
  { name: 'flag', category: '基础', path: 'M5 3v18H3V3h2zm2 0h11l-2 4 2 4H7V3z' },
  { name: 'eye', category: '基础', path: 'M12 5C6 5 2 12 2 12s4 7 10 7 10-7 10-7-4-7-10-7zm0 11a4 4 0 1 1 0-8 4 4 0 0 1 0 8z' },
  { name: 'edit', category: '基础', path: 'M3 17.25V21h3.75L17.8 9.94l-3.75-3.75L3 17.25zM20.7 7.04a1 1 0 0 0 0-1.41l-2.34-2.34a1 1 0 0 0-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z' },
  { name: 'trash', category: '基础', path: 'M6 7h12l-1 14H7L6 7zm3-3h6l1 2H8l1-2zM4 6h16v2H4V6z' },
  { name: 'plus', category: '基础', path: 'M11 5h2v6h6v2h-6v6h-2v-6H5v-2h6V5z' },
  { name: 'minus', category: '基础', path: 'M5 11h14v2H5v-2z' },
  { name: 'check', category: '基础', path: 'M9 16.2L4.8 12l-1.4 1.4L9 19 21 7l-1.4-1.4z' },
  { name: 'close', category: '基础', path: 'M18.3 5.7L12 12l6.3 6.3-1.4 1.4L10.6 13.4 4.3 19.7 2.9 18.3 9.2 12 2.9 5.7 4.3 4.3l6.3 6.3 6.3-6.3z' },
  { name: 'info', category: '基础', path: 'M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z' },
  { name: 'warning', category: '基础', path: 'M12 2L1 21h22L12 2zm1 15h-2v-2h2v2zm0-4h-2v-4h2v4z' },
  { name: 'question', category: '基础', path: 'M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20zm0 16a1.2 1.2 0 1 1 0-2.4 1.2 1.2 0 0 1 0 2.4zm1.4-5.6c-.6.4-.8.7-.8 1.4v.6h-1.2v-.6c0-1.1.5-1.8 1.3-2.4.7-.5 1-.8 1-1.4 0-.8-.7-1.4-1.5-1.4-.7 0-1.3.4-1.6 1l-1.1-.7C9.4 7.8 10.6 7 12 7c2 0 3.5 1.3 3.5 3.1 0 1.4-.8 2.2-1.6 2.8z' },

  // 文件
  { name: 'file', category: '文件', path: 'M6 2h8l6 6v14H6V2zm7 1.5V9h5.5L13 3.5z' },
  { name: 'folder', category: '文件', path: 'M3 5h6l2 2h10v12H3V5z' },
  { name: 'paperclip', category: '文件', path: 'M16.5 6.5l-8 8a3 3 0 0 0 4.2 4.2l8-8a5 5 0 0 0-7.1-7.1l-8 8a7 7 0 0 0 9.9 9.9l7.1-7.1-1.4-1.4-7.1 7.1a5 5 0 0 1-7.1-7.1l8-8a3 3 0 0 1 4.2 4.2l-8 8a1 1 0 0 1-1.4-1.4l8-8 1.4 1.4z' },
  { name: 'printer', category: '文件', path: 'M7 3h10v4H7V3zm-4 5h18a1 1 0 0 1 1 1v8h-4v4H6v-4H2v-8a1 1 0 0 1 1-1zm3 3v6h8v-6H6z' },
  { name: 'download', category: '文件', path: 'M12 3v10l4-4 1.4 1.4L12 15.8 6.6 10.4 8 9l4 4V3h0zM4 19h16v2H4v-2z' },
  { name: 'upload', category: '文件', path: 'M12 21V11l-4 4-1.4-1.4L12 8.2l5.4 5.4L16 15l-4-4v10h0zM4 3h16v2H4V3z' },
  { name: 'cloud', category: '文件', path: 'M7 18a4 4 0 0 1-.5-7.97 5.5 5.5 0 0 1 10.6-1.2A4.5 4.5 0 0 1 17 18H7z' },
  { name: 'cloud-upload', category: '文件', path: 'M7 18a4 4 0 0 1-.5-7.97 5.5 5.5 0 0 1 10.6-1.2A4.5 4.5 0 0 1 17 18h-2v-2h2v-1h-2v-2l-3 3-3-3v2h-2v1h2v2h-2zm-2-4h2v-2h2v2h2v-2h2v2h2v-1h-2V9H7v5H5v0z' },

  // 图表
  { name: 'bar', category: '图表', path: 'M4 20h16v2H4v-2zm2-2h3v-7H6v7zm5 0h3V6h-3v12zm5 0h3v-10h-3v10z' },
  { name: 'line', category: '图表', path: 'M3 19l6-7 4 4 7-9v3l-6 8-4-4-5 6H3z' },
  { name: 'pie', category: '图表', path: 'M12 2a10 10 0 1 0 10 10h-10V2zm0 2v8h8a8 8 0 0 1-8-8z' },
  { name: 'trend-up', category: '图表', path: 'M3 17l6-6 4 4 8-8v6h-2V8.4l-6 6-4-4-4 4-2-1z' },
  { name: 'trend-down', category: '图表', path: 'M3 7l6 6 4-4 8 8v-6h-2v2.6l-6-6-4 4-4-4-2 1z' },
  { name: 'dashboard', category: '图表', path: 'M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18zm0 2v8h8a8 8 0 0 1-8-8z' },

  // 媒体
  { name: 'image', category: '媒体', path: 'M4 4h16v16H4V4zm2 2v8l4-4 3 3 4-5 1 1v7H6V6zm2 2a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3z' },
  { name: 'video', category: '媒体', path: 'M3 5h14a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1zm15 3l4 4-4 4V8z' },
  { name: 'music', category: '媒体', path: 'M9 17V5l11-2v12a3 3 0 1 1-2-2.8V6.3L11 7.5V17a3 3 0 1 1-2-2.8z' },
  { name: 'camera', category: '媒体', path: 'M9 4l-1.5 2H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h16a1 1 0 0 0 1-1V7a1 1 0 0 0-1-1h-3.5L15 4H9zm3 4a5 5 0 1 1 0 10 5 5 0 0 1 0-10zm0 2a3 3 0 1 0 0 6 3 3 0 0 0 0-6z' },
  { name: 'play', category: '媒体', path: 'M8 5v14l11-7L8 5z' },
  { name: 'pause', category: '媒体', path: 'M7 5h3v14H7V5zm7 0h3v14h-3V5z' },

  // 通讯
  { name: 'mail', category: '通讯', path: 'M3 5h18a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1zm9 7L4 7v.01L12 12l8-5V7l-8 5z' },
  { name: 'phone', category: '通讯', path: 'M6 3a2 2 0 0 1 2 2c0 1.2.4 2.4 1 3.4.4.6 1 1 1.6 1 .7-.6 1.4-1 2.2-1-.2-1 .1-2 .8-2.7A14 14 0 0 0 6 3zm0 2.5l.3.5c.5.9.9 1.9 1.2 3-.3 0-.6.2-.9.5l-.4.4-.4-.3-.4-.5A14 14 0 0 0 6 5.5z' },
  { name: 'message', category: '通讯', path: 'M4 4h16a1 1 0 0 1 1 1v11a1 1 0 0 1-1 1H9l-5 4v-4H4a1 1 0 0 1-1-1V5a1 1 0 0 1 1-1z' },
  { name: 'send', category: '通讯', path: 'M3 3l18 9-18 9 4-9-4-9zm4 9l-2 0 0 0z' },
  { name: 'share', category: '通讯', path: 'M18 8a3 3 0 1 0-2.8-4L8.9 7.2a3 3 0 1 0 0 5.6l6.3 3.2A3 3 0 1 0 18 16a3 3 0 0 0-2.8 2l-6.3-3.2a3 3 0 1 0 0-5.6L15.2 8A3 3 0 0 0 18 8z' },

  // 商业
  { name: 'cart', category: '商业', path: 'M3 3h2l2.5 12h11L21 7H6.2l-.5-2H3v2zm4 16a2 2 0 1 0 0 4 2 2 0 0 0 0-4zm10 0a2 2 0 1 0 0 4 2 2 0 0 0 0-4z' },
  { name: 'card', category: '商业', path: 'M3 5h18a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1zm2 4v2h6V9H5z' },
  { name: 'money', category: '商业', path: 'M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20zm0 3a3 3 0 1 1 0 6 3 3 0 0 1 0-6zm0 9a4 4 0 0 1-3-1.4c.6-.4 1-1 1-1.6h4c0 .6.4 1.2 1 1.6A4 4 0 0 1 12 14z' },
  { name: 'tag', category: '商业', path: 'M3 3h10l8 8-8 8-8-8V3zm3 3a2 2 0 1 0 0 4 2 2 0 0 0 0-4z' },
  { name: 'gift', category: '商业', path: 'M12 8a3 3 0 0 1 3 3v9h-6v-9a3 3 0 0 1 3-3zm-7 3h5v9H5V11zm14 0h-5v9h5V11zM12 8V4h0a3 3 0 0 1 0 4z' },
  { name: 'briefcase', category: '商业', path: 'M3 7h18a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V8a1 1 0 0 1 1-1zm4 0V5a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v2H7z' },

  // 系统
  { name: 'database', category: '系统', path: 'M12 3c4 0 8 1 8 3s-4 3-8 3-8-1-8-3 4-3 8-3zm8 5c0 1.5-3.6 2.7-8 2.7S4 9.5 4 8v5c0 1.5 3.6 2.7 8 2.7s8-1.2 8-2.7V8zm0 7c0 1.5-3.6 2.7-8 2.7S4 16.5 4 15v5c0 1.5 3.6 2.7 8 2.7s8-1.2 8-2.7v-5z' },
  { name: 'server', category: '系统', path: 'M4 3h16a1 1 0 0 1 1 1v6a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1zm1 3a1 1 0 1 0 0-2 1 1 0 0 0 0 2zm0 8a1 1 0 1 0 0-2 1 1 0 0 0 0 2zM3 14h18a1 1 0 0 1 1 1v6a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1v-6a1 1 0 0 1 1-1z' },
  { name: 'code', category: '系统', path: 'M9.4 7.4L4.8 12l4.6 4.6 1.4-1.4L7.6 12l3.2-3.2-1.4-1.4zm5.2 0L13.2 9.4 16.4 12l-3.2 3.2 1.4 1.4L19.2 12l-4.6-4.6z' },
  { name: 'terminal', category: '系统', path: 'M4 3h16a1 1 0 0 1 1 1v16a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1zm3 5l3 3-3 3 1.4 1.4L13 12l-4.6-4.6L7 8zM14 16h4v2h-4v-2z' },
  { name: 'bug', category: '系统', path: 'M12 4a2 2 0 0 1 2 2v1.1a5 5 0 0 1 3 4.6V15h2v2h-2v1a2 2 0 0 1-2 2h-1v-2h-2v2H9a2 2 0 0 1-2-2v-1H5v-2h2v-3.3A5 5 0 0 1 8 7.1V6a2 2 0 0 1 2-2h2zm-3 5h6a3 3 0 0 1-6 0z' },
  { name: 'shield', category: '系统', path: 'M12 2l8 3v6c0 5-3.5 8.5-8 11-4.5-2.5-8-6-8-11V5l8-3zm0 2.2L6 6.3V11c0 3.7 2.5 6.4 6 8.4 3.5-2 6-4.7 6-8.4V6.3l-6-2.1z' },
  { name: 'lock', category: '系统', path: 'M6 10V8a6 6 0 0 1 12 0v2h1a1 1 0 0 1 1 1v9a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1v-9a1 1 0 0 1 1-1h1zm2 0h8V8a4 4 0 0 0-8 0v2zm4 4a2 2 0 0 1 1 3.7V18h-2v-2.3A2 2 0 0 1 12 14z' },
  { name: 'key', category: '系统', path: 'M14 2a6 6 0 0 0-5.6 8.2L2 16.6V22h5.4l3.4-3.4V17h2v-2.4l1.2-1.2A6 6 0 1 0 14 2zm2.5 3.5a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3z' },

  // 位置时间
  { name: 'map', category: '位置', path: 'M12 2a7 7 0 0 0-7 7c0 5 7 13 7 13s7-8 7-13a7 7 0 0 0-7-7zm0 9.5A2.5 2.5 0 1 1 12 6a2.5 2.5 0 0 1 0 5.5z' },
  { name: 'pin', category: '位置', path: 'M12 2a7 7 0 0 0-7 7c0 5 7 13 7 13s7-8 7-13a7 7 0 0 0-7-7zm0 9.5A2.5 2.5 0 1 1 12 6a2.5 2.5 0 0 1 0 5.5z' },
  { name: 'location', category: '位置', path: 'M12 2a7 7 0 0 0-7 7c0 5 7 13 7 13s7-8 7-13a7 7 0 0 0-7-7zm0 9.5A2.5 2.5 0 1 1 12 6a2.5 2.5 0 0 1 0 5.5z' },
  { name: 'globe', category: '位置', path: 'M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20zm0 2c1.5 0 3 2.2 3.4 5H8.6C9 6.2 10.5 4 12 4zM4.3 11h3.1c.1-1.9.5-3.6 1-5A8 8 0 0 0 4.3 11zm0 2a8 8 0 0 0 4.1 5c-.5-1.4-.9-3.1-1-5H4.3zm3.3 7c-.5-1.4-.9-3.1-1-5h3.1c.1 2.6.6 4.7 1 5H7.6zm4.4 0c.4-.3.9-2.4 1-5h3.1c-.1 1.9-.5 3.6-1 5H12zm4.1-7c-.1 1.9-.5 3.6-1 5a8 8 0 0 0 4.1-5h-3.1zm1-2H13.4c-.4-2.8-.9-4.7-1-5a8 8 0 0 1 4.1 5z' },
  { name: 'calendar', category: '位置', path: 'M5 3h14a1 1 0 0 1 1 1v16a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1zm1 2v2h12V5H6zm0 4v10h12V9H6zm2 2h2v2H8v-2z' },
  { name: 'clock', category: '位置', path: 'M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20zm0 2a8 8 0 1 1 0 16 8 8 0 0 1 0-16zm1 3v5l4 2-.8 1.6L11 14V7h2z' },

  // 品牌/Misc
  { name: 'github', category: '品牌', path: 'M12 2a10 10 0 0 0-3.2 19.5c.5.1.7-.2.7-.5v-1.7c-2.8.6-3.4-1.4-3.4-1.4-.4-1.2-1.1-1.5-1.1-1.5-.9-.6.1-.6.1-.6 1 .1 1.5 1 1.5 1 .9 1.5 2.3 1.1 2.9.8.1-.6.3-1.1.6-1.4-2.2-.2-4.6-1.1-4.6-5 0-1.1.4-2 1-2.7-.1-.3-.4-1.3.1-2.7 0 0 .8-.3 2.7 1a9.4 9.4 0 0 1 5 0c1.9-1.3 2.7-1 2.7-1 .5 1.4.2 2.4.1 2.7.6.7 1 1.6 1 2.7 0 3.9-2.4 4.8-4.6 5 .3.3.6.9.6 1.9v2.8c0 .3.2.6.7.5A10 10 0 0 0 12 2z' },
  { name: 'apple', category: '品牌', path: 'M16.5 12.6c0-2 1.6-3 1.7-3-1-1.4-2.4-1.6-2.9-1.6-1.2-.1-2.4.7-3 .7-.6 0-1.6-.7-2.6-.7-1.3 0-2.6.8-3.3 2-1.4 2.4-.4 6 1 8 .7 1 1.4 2 2.5 2 1 0 1.3-.6 2.5-.6 1.2 0 1.5.6 2.5.6 1 0 1.7-1 2.4-2 .5-.7.7-1.4.7-1.4-1.6-.6-1.5-2.3-1.5-2.3zM14.4 5.8c.5-.7.9-1.6.8-2.5-.8 0-1.7.5-2.3 1.2-.5.6-.9 1.5-.8 2.4.9.1 1.7-.4 2.3-1.1z' },
  { name: 'android', category: '品牌', path: 'M5 9a3 3 0 0 1 3-3h8a3 3 0 0 1 3 3v9a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1V9zm2.5-5l1 1h7l1-1 1 1-2 2h-7l-2-2 1-1zM8 12v4h2v-4H8zm6 0v4h2v-4h-2z' },
  { name: 'windows', category: '品牌', path: 'M3 3h8v8H3V3zm9 0h9v8h-9V3zM3 12h8v9H3v-9zm9 0h9v9h-9v-9z' },
  { name: 'wifi', category: '品牌', path: 'M2 8.5a14 14 0 0 1 20 0l-1.4 1.4a12 12 0 0 0-17.2 0L2 8.5zm3 3a9 9 0 0 1 14 0l-1.5 1.5a7 7 0 0 0-11 0L5 11.5zm3 3a4 4 0 0 1 8 0l-1.5 1.5a2 2 0 0 0-5 0L8 15zM12 18a1 1 0 1 0 0 2 1 1 0 0 0 0-2z' },
  { name: 'battery', category: '品牌', path: 'M3 8h16a2 2 0 0 1 2 2v4a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2v-4a2 2 0 0 1 2-2zm18 3h2v4h-2v-4zM5 10v4h12v-4H5z' },
  { name: 'sun', category: '品牌', path: 'M12 7a5 5 0 1 0 0 10 5 5 0 0 0 0-10zm0-5h0v3h0V2zm0 17v3h0v-3zM4.2 4.2l2.1 2.1-1.4 1.4L2.8 5.6 4.2 4.2zm14.6 14.6l2.1 2.1-1.4 1.4-2.1-2.1 1.4-1.4zM2 12h3v0H2zm17 0h3v0h-3zM4.2 19.8l1.4-1.4 2.1 2.1-1.4 1.4-2.1-2.1zm14.6-14.6l1.4-1.4 2.1 2.1-1.4 1.4-2.1-2.1z' },
  { name: 'moon', category: '品牌', path: 'M12 3a9 9 0 1 0 9 9 7 7 0 0 1-9-9z' },
  { name: 'zap', category: '品牌', path: 'M13 2L3 14h7l-1 8 10-12h-7l1-8z' },
  { name: 'fire', category: '品牌', path: 'M12 2c1 3-1 4-2 6-1 2 0 3 1 4 0-2 1-3 2-3s2 1 2 3c0 3-2 5-5 5s-5-2-5-5c0-4 4-6 5-10 1 1 1 2 2 0z' },
];

export const iconCategories = Array.from(new Set(localIcons.map((i) => i.category)));

// 模拟"数百个"演示：通过组合变体生成扩展图标（仅用于专栏展示演示）
export function expandLocalIcons(times = 4): LocalIconDef[] {
  const result: LocalIconDef[] = [];
  let idx = 0;
  for (let t = 0; t < times; t += 1) {
    localIcons.forEach((ic) => {
      idx += 1;
      result.push({
        name: `${ic.name}${t === 0 ? '' : `-v${t}`}`,
        category: ic.category,
        path: ic.path,
      });
    });
  }
  return result;
}
