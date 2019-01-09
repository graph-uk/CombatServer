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
set PATH=%PATH%;C:\cygwin64\bin
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


call npm install
call npm run build

pushd %GOPATH%\src\github.com\graph-uk\combat
go build -ldflags="-s -w"
upx --brute combat.exe
popd

pushd %GOPATH%\src\github.com\graph-uk\combat-client
go build -ldflags="-s -w"
upx --brute combat-client.exe
popd

pushd %GOPATH%\src\github.com\graph-uk\combat-worker
go build -ldflags="-s -w"
upx --brute combat-worker.exe
popd


mkdir assets\_\dist
mkdir assets\_\dist\win64\
mkdir assets\_\dist\linux64\

copy %GOPATH%\src\github.com\graph-uk\combat\combat.exe assets\_\dist\win64\
copy %GOPATH%\src\github.com\graph-uk\combat-client\combat-client.exe assets\_\dist\win64\
copy %GOPATH%\src\github.com\graph-uk\combat-worker\combat-worker.exe assets\_\dist\win64\

docker rm -f combat-builder
docker run --rm --name combat-builder -v %GOPATH%\src:/go/src golang:1.9.2 bash -c "cd /go/src/github.com/graph-uk/combat-server && apt-get update && apt-get -y install upx dos2unix && dos2unix ./*.sh &&./build_linux_binaries.sh"

copy %GOPATH%\src\github.com\graph-uk\combat\combat assets\_\dist\linux64\
copy %GOPATH%\src\github.com\graph-uk\combat-client\combat-client assets\_\dist\linux64\
copy %GOPATH%\src\github.com\graph-uk\combat-worker\combat-worker assets\_\dist\linux64\

pushd %GOPATH%\src\github.com\graph-uk\combat-server
packr build -ldflags="-s -w"
upx --brute combat-server.exe
popd
