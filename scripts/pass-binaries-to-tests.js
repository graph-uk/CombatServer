const t = new(require('./tools/tools.js'))

const targetFolder = process.env['REPO']+'/int-tests/src/Tests_shared/target-app-binaries'

t.rmdir(targetFolder)
t.mkdir(targetFolder)


if (process.platform === "win32"){
	t.mv(process.env['REPO']+'/src/malibu/malibu.exe', targetFolder+'/malibu.exe')
	t.mv(process.env['REPO']+'/src/malibu-client/malibu-client.exe', targetFolder+'/malibu-client.exe')
	t.mv(process.env['REPO']+'/src/malibu-server/malibu-server.exe', targetFolder+'/malibu-server.exe')
	t.mv(process.env['REPO']+'/src/malibu-worker/malibu-worker.exe', targetFolder+'/malibu-worker.exe')
}else{
	t.mv(process.env['REPO']+'/src/malibu/malibu', targetFolder+'/malibu')
	t.mv(process.env['REPO']+'/src/malibu-client/malibu-client', targetFolder+'/malibu-client')
	t.mv(process.env['REPO']+'/src/malibu-server/malibu-server', targetFolder+'/malibu-server')
	t.mv(process.env['REPO']+'/src/malibu-worker/malibu-worker', targetFolder+'/malibu-worker')	
}
