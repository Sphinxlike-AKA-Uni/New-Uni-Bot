@echo off
go build Uni.go
if %errorlevel% EQU 0 "Uni.exe" -config ../UniConfig.inf
Pause