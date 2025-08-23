package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetFileHash 计算文件的MD5和SHA256哈希值
func GetFileHash(file multipart.File) (string, string, error) {
	// 重置文件指针到开始位置
	file.Seek(0, 0)
	
	// 创建哈希计算器
	md5Hash := md5.New()
	sha256Hash := sha256.New()
	
	// 使用MultiWriter同时写入两个哈希计算器
	multiWriter := io.MultiWriter(md5Hash, sha256Hash)
	
	// 复制文件内容到哈希计算器
	_, err := io.Copy(multiWriter, file)
	if err != nil {
		return "", "", err
	}
	
	// 重置文件指针到开始位置
	file.Seek(0, 0)
	
	return fmt.Sprintf("%x", md5Hash.Sum(nil)), fmt.Sprintf("%x", sha256Hash.Sum(nil)), nil
}

// GenerateFileName 生成唯一的文件名
func GenerateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	nameWithoutExt := strings.TrimSuffix(originalName, ext)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%d%s", nameWithoutExt, timestamp, ext)
}

// GetMimeType 根据文件扩展名获取MIME类型
func GetMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".txt":  "text/plain",
		".csv":  "text/csv",
		".mp4":  "video/mp4",
		".avi":  "video/avi",
		".mov":  "video/quicktime",
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".zip":  "application/zip",
		".rar":  "application/x-rar-compressed",
		".7z":   "application/x-7z-compressed",
	}
	
	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}
	
	return "application/octet-stream"
}

// GetFileType 根据MIME类型获取文件类型分类
func GetFileType(mimeType string) string {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return "image"
	case strings.HasPrefix(mimeType, "video/"):
		return "video"
	case strings.HasPrefix(mimeType, "audio/"):
		return "audio"
	case strings.HasPrefix(mimeType, "text/"):
		return "text"
	case mimeType == "application/pdf":
		return "pdf"
	case strings.Contains(mimeType, "document") || strings.Contains(mimeType, "word"):
		return "document"
	case strings.Contains(mimeType, "sheet") || strings.Contains(mimeType, "excel"):
		return "spreadsheet"
	case strings.Contains(mimeType, "presentation") || strings.Contains(mimeType, "powerpoint"):
		return "presentation"
	case strings.Contains(mimeType, "zip") || strings.Contains(mimeType, "compressed"):
		return "archive"
	default:
		return "other"
	}
}

// EnsureDir 确保目录存在，如果不存在则创建
func EnsureDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// FormatFileSize 格式化文件大小
func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}