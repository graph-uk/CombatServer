const t = new(require('./tools/tools.js'))

const targetFolder = process.env['REPO']+'/int-tests/src/Tests_shared/target-app-binaries'

t.rmdir(targetFolder)
t.mkdir(targetFolder)
t.mv(process.env['REPO']+'/src/malibu/malibu.exe', targetFolder+'/malibu.exe')
t.mv(process.env['REPO']+'/src/malibu-client/malibu-client.exe', targetFolder+'/malibu-client.exe')
t.mv(process.env['REPO']+'/src/malibu-server/malibu-server.exe', targetFolder+'/malibu-server.exe')
t.mv(process.env['REPO']+'/src/malibu-worker/malibu-worker.exe', targetFolder+'/malibu-worker.exe')