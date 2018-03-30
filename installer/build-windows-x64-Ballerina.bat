@echo off

set BALLERINA_VERSION=0.970.0-alpha1
set BALZIP=../distribution/zip/ballerina/target/ballerina-%BALLERINA_VERSION%.zip
set BALDIST=ballerina-%BALLERINA_VERSION%
set BALPOS=windows
set BALPARCH=x64

for /f %%x in ('wmic path win32_utctime get /format:list ^| findstr "="') do set %%x
set UTC_TIME=%Year%-%Month%-%Day% %Hour%:%Minute%:%Second% UTC

del msi\%BALDIST%-%BALPOS%-%BALPARCH%.msi /s /q >nul 2>&1
rmdir ballerina-%BALLERINA_VERSION% /s /q >nul 2>&1

echo %BALDIST% build started at '%UTC_TIME%'
go run -ldflags "-X main.balVersion=%BALLERINA_VERSION% -X main.balPOS=%BALPOS% -X main.balPArch=%BALPARCH% -X main.balZip=%BALZIP% -X main.balDist=%BALDIST%" release.go
echo %BALDIST% build completed at '%UTC_TIME%'