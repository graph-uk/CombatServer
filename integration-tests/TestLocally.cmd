set GOPATH=%cd%\..\..\..\..\..\
set GOROOT=%GOPATH%\combat-dev-utils\combat-dev-go
set PATH=%PATH%;%GOROOT%\bin

go run test.go
pause