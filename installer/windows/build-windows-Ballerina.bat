@echo off

:argumentLoop
IF NOT "%1"=="" (
    IF "%1"=="-dist" (
        SET DIST=%2
        SHIFT
    )
    IF "%1"=="-version" (
        SET BALLERINA_VERSION=%2
        SHIFT
    )
	SHIFT
	goto argumentLoop
) 


IF "%DIST%"==""  (
	echo The syntax of the command is incorrect. Missing argument dist.
	goto EOF
)

IF NOT "%DIST%"=="all" IF NOT "%DIST%"=="BallerinaToolsWin64" IF NOT "%DIST%"=="BallerinaToolsWin586" IF NOT "%DIST%"=="BallerinaWin64" IF NOT "%DIST%"=="BallerinaWin586" (
	echo The syntax of the command is incorrect. Possible arguments for dist - all, BallerinaToolsWin64, BallerinaToolsWin586, BallerinaWin64, BallerinaWin586.
	echo Ex: -dist BallerinaToolsWin64
	goto EOF
)

IF "%BALLERINA_VERSION%"==""  (
	echo The syntax of the command is incorrect. Missing argument version.
	goto EOF
)

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

IF "%DIST%"=="all" (
	call :createBallerinaToolsWin64Installer
	call :createBallerinaToolsWin586Installer
	call :createBallerinaWin64Installer
	call :createBallerinaWin586Installer
) ELSE (
	call :create%DIST%Installer
)

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
rem jar -xf %BALZIP%
rmdir ballerina-%BALLERINA_VERSION% /s /q >nul 2>&1
rmdir ballerina-tools-%BALLERINA_VERSION% /s /q >nul 2>&1
rmdir target\installer-resources /s /q >nul 2>&1
powershell.exe -nologo -noprofile -command "& { Add-Type -A 'System.IO.Compression.FileSystem'; [IO.Compression.ZipFile]::ExtractToDirectory('%BALZIP%', '.'); }"

echo %BALDIST% build started at '%UTC_TIME%' for %BALPOS% %BALPARCH%
%WIXDIST%\heat dir %BALDIST% -nologo -gg -g1 -srd -sfrag -sreg -cg AppFiles -template fragment -dr INSTALLDIR -var var.SourceDir -out target\installer-resources\AppFiles.wxs
%WIXDIST%\candle -nologo -dbalVersion=%BALLERINA_VERSION% -dWixbalVersion=1.0.0.0 -dArch=%INSTALLERPARCH% -dSourceDir=%BALDIST% -out target\installer-resources\ -ext WixUtilExtension resources\installer.wxs target\installer-resources\AppFiles.wxs
%WIXDIST%\light -nologo -dcl:high -sice:ICE60 -ext WixUIExtension -ext WixUtilExtension -loc resources\en-us.wxl target\installer-resources\AppFiles.wixobj target\installer-resources\installer.wixobj -o target\msi\%BALDIST%-%BALPOS%-%BALPARCH%.msi
echo %BALDIST% build completed at '%UTC_TIME%' for %BALPOS% %BALPARCH%
echo.
goto EOF

:EOF
