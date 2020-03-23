rem Run "installCombatDevEnv.cmd" as admin before

::change dir to the BAT directory
cd /D %~dp0

@echo off
set GOPATH=%cd%\..\..\..\..\
set GOROOT=%GOPATH%\node_modules\combat-dev-go
set PATH=%PATH%;%GOROOT%\bin
set PATH=%PATH%;%GOPATH%\bin
set PATH=%PATH%;%GOPATH%\node_modules\combat-dev-mingw64\mingw64\bin
set PATH=%PATH%;%GOPATH%\node_modules\Nodejs.Redist.x64\tools
@echo on

del /F /S /Q assets\_
del /F /S /Q combat-server.exe

call npm install
call npm run build

pushd %GOPATH%\src\github.com\graph-uk\combat
start go build
popd

pushd %GOPATH%\src\github.com\graph-uk\combat-client
start go build
popd

pushd %GOPATH%\src\github.com\graph-uk\combat-worker
start go build
popd

pushd %GOPATH%\src\github.com\graph-uk\combat-server
packr build
popd

pushd integration-tests
TestLocally.cmd
popd