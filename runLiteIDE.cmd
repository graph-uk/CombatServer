rem Run "installCombatDevEnv.cmd" as admin before

set GOPATH=%cd%\..\..\..\..\
set GOROOT=%GOPATH%\combat-dev-utils\combat-dev-go
set PATH=%PATH%;%GOROOT%\bin
set PATH=%PATH%;%GOPATH%\combat-dev-utils\combat-dev-mingw64\mingw64\bin
start %GOPATH%\combat-dev-utils\combat-dev-liteide\liteide\bin\liteide.exe