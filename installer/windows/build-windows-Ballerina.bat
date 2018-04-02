@echo off

set BALLERINA_VERSION=0.970.0-alpha1-SNAPSHOT
set BALPOS=windows
set WIXDIST=resources\wix

rem set BALZIP=..\..\distribution\zip\ballerina\target\ballerina-%BALLERINA_VERSION%.zip
rem set BALDIST=ballerina-%BALLERINA_VERSION%
rem set BALPARCH=x64

for /f %%x in ('wmic path win32_utctime get /format:list ^| findstr "="') do set %%x
set UTC_TIME=%Year%-%Month%-%Day% %Hour%:%Minute%:%Second% UTC

del msi\%BALDIST%-%BALPOS%-%BALPARCH%.msi /s /q >nul 2>&1
rmdir ballerina-%BALLERINA_VERSION% /s /q >nul 2>&1
rmdir ballerina-tools-%BALLERINA_VERSION% /s /q >nul 2>&1
rmdir target /s /q >nul 2>&1

call :createBallerinaToolsWin64Installer
call :createBallerinaToolsWin586Installer
call :createBallerinaWin64Installer
call :createBallerinaWin586Installer

goto EOF


:createBallerinaToolsWin64Installer
set BALZIP=..\..\distribution\zip\ballerina-tools\target\ballerina-tools-%BALLERINA_VERSION%.zip
set BALDIST=ballerina-tools-%BALLERINA_VERSION%
set BALPARCH=x64
set INSTALLERPARCH=amd64
call :createInstaller
goto EOF

:createBallerinaToolsWin586Installer
set BALZIP=..\..\distribution\zip\ballerina-tools\target\ballerina-tools-%BALLERINA_VERSION%.zip
set BALDIST=ballerina-tools-%BALLERINA_VERSION%
set BALPARCH=i586
set INSTALLERPARCH=386
call :createInstaller
goto EOF

:createBallerinaWin64Installer
set BALZIP=..\..\distribution\zip\ballerina\target\ballerina-%BALLERINA_VERSION%.zip
set BALDIST=ballerina-%BALLERINA_VERSION%
set BALPARCH=x64
set INSTALLERPARCH=amd64
call :createInstaller
goto EOF

:createBallerinaWin586Installer
set BALZIP=..\..\distribution\zip\ballerina\target\ballerina-%BALLERINA_VERSION%.zip
set BALDIST=ballerina-%BALLERINA_VERSION%
set BALPARCH=i586
set INSTALLERPARCH=386
call :createInstaller
goto EOF

:createInstaller
jar -xf %BALZIP%

echo %BALDIST% build started at '%UTC_TIME%' for %BALPOS% %BALPARCH%
%WIXDIST%\heat dir %BALDIST% -nologo -gg -g1 -srd -sfrag -sreg -cg AppFiles -template fragment -dr INSTALLDIR -var var.SourceDir -out target\AppFiles.wxs
%WIXDIST%\candle -nologo -dbalVersion=%BALLERINA_VERSION% -dWixbalVersion=1.0.0.0 -dArch=%INSTALLERPARCH% -dSourceDir=%BALDIST% -out target\ -ext WixUtilExtension resources\installer.wxs target\AppFiles.wxs
%WIXDIST%\light -nologo -dcl:high -sice:ICE60 -ext WixUIExtension -ext WixUtilExtension -loc resources\en-us.wxl target\AppFiles.wixobj target\installer.wixobj -o target\msi\%BALDIST%-%BALPOS%-%BALPARCH%.msi
echo %BALDIST% build completed at '%UTC_TIME%' for %BALPOS% %BALPARCH%
echo.
goto EOF

:EOF
