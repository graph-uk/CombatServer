const t = new(require('./tools/tools.js'));
const r = process.env.REPO

t.cmd('packr build', {cwd: r+'/src/malibu-server'})
