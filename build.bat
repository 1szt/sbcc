@echo off
:: 切换命令行编码为 UTF-8，解决中文显示问题
chcp 65001 >nul
SETLOCAL
REM ==========================================
REM Modbus 总线系统多平台自动构建脚本 (已修复 & 符号)
REM ==========================================

REM 设置参数
SET APP_NAME=sbcc
SET ENTRY_POINT=./modbus/main/run.go
SET BUILD_DIR=build

echo [1/4] 正在清理旧的构建任务...
if exist %BUILD_DIR% rd /s /q %BUILD_DIR%
mkdir %BUILD_DIR%

echo [2/4] 正在构建 Windows 版本 (64位)...
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags="-s -w" -o %BUILD_DIR%/windows/%APP_NAME%.exe %ENTRY_POINT%
if %errorlevel% neq 0 (echo Windows 构建失败! && pause && exit /b)

echo [3/4] 正在构建 Linux 版本 (64位)...
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags="-s -w" -o %BUILD_DIR%/linux/%APP_NAME% %ENTRY_POINT%
if %errorlevel% neq 0 (echo Linux 构建失败! && pause && exit /b)

:: 关键修正点：使用 ^ 符号转义 &
echo [4/4] 正在构建 Mac 版本 (Intel ^& M 系列)...
SET GOOS=darwin
SET GOARCH=amd64
go build -ldflags="-s -w" -o %BUILD_DIR%/mac_intel/%APP_NAME% %ENTRY_POINT%
SET GOARCH=arm64
go build -ldflags="-s -w" -o %BUILD_DIR%/mac_m1/%APP_NAME% %ENTRY_POINT%

echo.
echo ==========================================
echo 构建完成！产物已存放在 /%BUILD_DIR% 文件夹。
echo ==========================================
pause