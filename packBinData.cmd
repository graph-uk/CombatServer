set PATH=%PATH%;%cd%\vendor\github.com\jteeuwen\go-bindata\go-bindata
set PATH=%PATH%;%cd%\vendor\github.com\elazarl\go-bindata-assetfs\go-bindata-assetfs

:: go-bindata-assetfs.exe bindata/...
go-bindata-assetfs.exe bindata/... -o bindata_assetfs.go