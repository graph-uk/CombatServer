const t = new(require('./tools/tools.js'));
const r = process.env.REPO

goos = process.env.GOOS
goarch = process.env.GOARCH


process.env.GOARCH = 'amd64'

process.env.GOOS = 'linux'
t.cmd('go build', {cwd: r+'/src/malibu'})
process.env.GOOS = 'windows'
t.cmd('go build', {cwd: r+'/src/malibu'})

process.env.GOOS = 'linux'
t.cmd('go build', {cwd: r+'/src/malibu-client'})
process.env.GOOS = 'windows'
t.cmd('go build', {cwd: r+'/src/malibu-client'})

process.env.GOOS = 'linux'
t.cmd('go build', {cwd: r+'/src/malibu-worker'})
process.env.GOOS = 'windows'
t.cmd('go build', {cwd: r+'/src/malibu-worker'})

process.chdir(r);
process.env.GOOS = goos
process.env.GOARCH = goarch
