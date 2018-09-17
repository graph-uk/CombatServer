rem Run "installCombatDevEnv.cmd" as admin before

::change dir to the BAT directory
cd /D %~dp0

@echo off
set GOPATH=%cd%\..\..\..\..\
set GOROOT=%GOPATH%\combat-dev-utils\combat-dev-go
set PATH=%PATH%;%GOROOT%\bin
set PATH=%PATH%;%GOPATH%\combat-dev-utils\combat-dev-mingw64\mingw64\bin
@echo on

echo "TESTLine"

pushd %GOPATH%\src\github.com\graph-uk\combat
start go build
popd

pushd %GOPATH%\src\github.com\graph-uk\combat-client
start go build
popd

pushd %GOPATH%\src\github.com\graph-uk\combat-worker
start go build
popd

call packBinData.cmd

pushd %GOPATH%\src\github.com\graph-uk\combat-server
go build
popd

pushd integration-tests
TestLocally.cmd
popd