@echo off
chcp 65001 >nul 2>&1
title 资料管理平台 - 启动脚本

echo ================================================
echo              资料管理平台启动脚本
echo ================================================
echo.

echo 检查环境依赖...
echo.

REM 检查Go环境
where go >nul 2>&1
if errorlevel 1 (
    echo [错误] 未找到Go环境，请确保已安装Go 1.21+
    echo.
    pause
    exit /b 1
)

REM 检查Node.js环境
where npm >nul 2>&1
if errorlevel 1 (
    echo [错误] 未找到Node.js环境，请确保已安装Node.js 18+
    echo.
    pause
    exit /b 1
)

echo [✓] 环境检查通过
echo.

echo 准备启动服务...
echo.
echo 后端服务将在 http://localhost:8081 启动
echo 前端服务将在 http://localhost:3000 启动
echo.
echo 提示: 使用 Ctrl+C 可以停止对应的服务
echo.

REM 启动后端服务（在新窗口中）
echo 正在启动后端服务...
start "资料管理平台-后端服务" /d "%~dp0backend" cmd /k "echo 启动Go后端服务器... & echo 当前目录: %cd% & echo. & echo 正在启动Go服务器 (端口: 8081)... & go run main.go"

REM 等待2秒让后端服务启动
timeout /t 2 /nobreak >nul

REM 启动前端服务（在新窗口中）
echo 正在启动前端服务...
start "资料管理平台-前端服务" /d "%~dp0frontend" cmd /k "echo 启动Vue前端开发服务器... & echo 当前目录: %cd% & echo. & echo 正在启动Vite开发服务器 (端口: 3000)... & npm run dev"

echo.
echo ================================================
echo 服务启动完成！
echo ================================================
echo.
echo 前端地址: http://localhost:3000
echo 后端地址: http://localhost:8081
echo.
echo 两个服务都在独立的窗口中运行
echo 关闭对应窗口或按 Ctrl+C 可停止服务
echo.

pause