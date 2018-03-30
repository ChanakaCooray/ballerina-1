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

jar -xf %BALZIP%

echo %BALDIST% build started at '%UTC_TIME%'
resources\wix\heat dir %BALDIST% -nologo -gg -g1 -srd -sfrag -sreg -cg AppFiles -template fragment -dr INSTALLDIR -var var.SourceDir -out target\AppFiles.wxs

resources\wix\candle -nologo -dbalVersion=%BALLERINA_VERSION% -dWixbalVersion=%BALLERINA_VERSION% -dArch=%BALPARCH% -dSourceDir=%BALDIST% -ext WixUtilExtension resources\installer.wxs

resources\wix\light -nologo -dc1:high -sice:ICE60 -ext WixUIExtension -ext WixUtilExtension -loc resources/en-us.wxl AppFiles.wixobj installer.wixobj -o target\msi\%BALDIST%-%BALPOS%-%BALPARCH%.msi

