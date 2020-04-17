const t = new(require('./tools/tools.js'))

//t.checkCmdAvailable('docker ps')
//"C:\Program Files (x86)\Google\Chrome\Application\chrome.exe" -incognito http://localhost/web/html/layout.html#/buckets
t.startCmdDetached('C:/Program Files (x86)/Google/Chrome/Application/chrome.exe', ['-incognito http://localhost:82/web/html/layout.html#/buckets'])
t.startCmdDetached('cmd',['/c start cmd /c boltdbweb.exe -p 82 -d '+process.env['REPO']+'/int-tests/src/Tests/twoSessions/server/malibu-base.db'])
//t.startCmdDetached('liteide', [process.env['REPO']+'/src/malibu-server/main.go'])