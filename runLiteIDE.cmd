rem Run "installCombatDevEnv.cmd" as admin before

set GOPATH=%cd%\..\..\..\..\
set GOROOT=%GOPATH%\node_modules\combat-dev-go
set PATH=%PATH%;%GOROOT%\bin
set PATH=%PATH%;%GOPATH%\bin
set PATH=%PATH%;%GOPATH%\node_modules\combat-dev-mingw64\mingw64\bin
start %GOPATH%\node_modules\combat-dev-liteide\liteide\bin\liteide.exe