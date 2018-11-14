rem Run the script as Admin
rem This script install globally:
rem 	-Chocolatey (with PS setups script https://chocolatey.org/install.ps1)
rem 	-Nuget (as Chocolatey package)
rem Locally, near "src", as nuget packets.
rem 	-go
rem 	-liteide
rem 	-mingw64
rem Locally, in project's GOPATH:
rem 	-combat
rem 	-combat-client
rem 	-combat-worker


:check Admin permissions
@echo off
echo Administrative permissions required. Detecting permissions...
net session >nul 2>&1
if %errorLevel% == 0 (
        echo Success: Administrative permissions confirmed.
) else (
        echo Failure: Please, re-run the script with Admin permissions!.
	pause >nul
	exit 1
)

::install chocolatey if command not found
choco.exe -v 2>NUL || @"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -InputFormat None -ExecutionPolicy Bypass -Command "iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))" && SET "PATH=%PATH%;%ALLUSERSPROFILE%\chocolatey\bin"

::install nuget if command not found
nuget.exe 2>NUL || choco install nuget.commandline

::change dir to the BAT directory
cd /D %~dp0/../../../..

if not exist "combat-dev-utils" mkdir "combat-dev-utils"
cd combat-dev-utils

::install combat-tests-dev-utils
nuget.exe install combat-dev-liteide -ExcludeVersion
nuget.exe install combat-dev-go -ExcludeVersion
nuget.exe install combat-dev-mingw64 -ExcludeVersion
nuget.exe install Nodejs.Redist.x64 -ExcludeVersion -Version 11.1.0

:: build bindata builder
cd /D %~dp0

set GOPATH=%cd%\..\..\..\..
set GOROOT=%GOPATH%\combat-dev-utils\combat-dev-go
set PATH=%PATH%;%GOROOT%\bin

pushd vendor\github.com\elazarl\go-bindata-assetfs\go-bindata-assetfs
go build
popd
pushd vendor\github.com\jteeuwen\go-bindata\go-bindata
go build
popd

:: get client and worker
go get github.com/graph-uk/combat
go get github.com/graph-uk/combat-client
go get github.com/graph-uk/combat-worker