@echo off
go build Uni.go
if %errorlevel% EQU 0 goto UniRun
goto End
:UniRun
"Uni.exe" -config ../UniConfig.inf
if %errorlevel% EQU 1 goto UniRun
rem ExitCode 1 means to restart UniBot
Pause