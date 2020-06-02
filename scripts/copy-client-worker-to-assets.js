const t = new(require('./tools/tools.js'))

const targetFolder = process.env['REPO']+'/src/malibu-server/assets/_/dist'

t.rmdir(targetFolder)
t.mkdir(targetFolder)
t.mkdir(targetFolder+'/client')
t.mkdir(targetFolder+'/worker')
t.mkdir(targetFolder+'/malibu')

t.cp(process.env['REPO']+'/src/malibu/malibu.exe', targetFolder+'/malibu/malibu.exe')
t.cp(process.env['REPO']+'/src/malibu/malibu', targetFolder+'/malibu/malibu')
t.cp(process.env['REPO']+'/src/malibu-client/malibu-client.exe', targetFolder+'/client/malibu-client.exe')
t.cp(process.env['REPO']+'/src/malibu-client/malibu-client', targetFolder+'/client/malibu-client')
t.cp(process.env['REPO']+'/src/malibu-worker/malibu-worker.exe', targetFolder+'/worker/malibu-worker.exe')
t.cp(process.env['REPO']+'/src/malibu-worker/malibu-worker', targetFolder+'/worker/malibu-worker')