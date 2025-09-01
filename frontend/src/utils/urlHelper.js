/**
 * 将相对URL转换为绝对URL
 * @param {string} url - 原始URL（可能是相对路径）
 * @returns {string} 绝对URL
 */
export function toAbsoluteUrl(url) {
  if (!url) return url
  
  // 如果已经是绝对URL，直接返回
  if (url.startsWith('http')) {
    return url
  }
  
  // 获取基础URL，如果环境变量未设置则使用默认值
  const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081'
  
  // 处理相对路径
  if (url.startsWith('/')) {
    // 以/开头的路径，直接拼接基础URL
    return baseUrl + url
  } else {
    // 不以/开头的路径，添加基础URL和/
    return baseUrl + '/' + url
  }
}