const t = new(require('./tools/tools.js'));
const r = process.env.REPO

t.cmd('malibu -HostName=localhost', {cwd: r+'/int-tests/src/Tests'})
