::Run "installCombatDevEnv.cmd" as admin before

::check docker available
docker ps 2>NUL || echo "Docker version returns non zero. Please check docker installed and available." && pause && exit 1

::change dir to the BAT directory
cd /D %~dp0

@echo off
set GOPATH=%cd%\..\..\..\..\
set GOROOT=%GOPATH%\combat-dev-utils\combat-dev-go
set PATH=%PATH%;%GOROOT%\bin
set PATH=%PATH%;%GOPATH%\bin
set PATH=%PATH%;%GOPATH%\combat-dev-utils\combat-dev-mingw64\mingw64\bin
set PATH=%PATH%;%GOPATH%\combat-dev-utils\Nodejs.Redist.x64\tools
set PATH=%PATH%;%GOPATH%\combat-dev-utils\combat-dev-upx
@echo on

:delete the old binaries
del /F /S /Q assets\_
del /F /S /Q combat-server.exe
del /F /S /Q combat-server
del /F /S /Q %GOPATH%\src\github.com\graph-uk\combat\combat*
del /F /S /Q %GOPATH%\src\github.com\graph-uk\combat-client\combat*
del /F /S /Q %GOPATH%\src\github.com\graph-uk\combat-worker\combat*




pushd %GOPATH%\src\github.com\graph-uk\combat
start go build
popd

pushd %GOPATH%\src\github.com\graph-uk\combat-client
start go build
popd

pushd %GOPATH%\src\github.com\graph-uk\combat-worker
start go build
popd

call npm install
call npm run build

pushd %GOPATH%\src\github.com\graph-uk\combat-server
packr build
popd

pushd integration-tests
TestLocally.cmd
popd
