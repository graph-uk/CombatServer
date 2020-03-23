set GOPATH=%cd%\..\..\..\..\..\
set GOROOT=%GOPATH%\node_modules\combat-dev-go
set PATH=%PATH%;%GOROOT%\bin

go run test.go
pause