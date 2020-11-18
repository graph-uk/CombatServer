const t = new(require('./tools/tools.js'));
const r = process.env.REPO

t.cmd('docker build -t malibu-worker .', {cwd: r+'/src/malibu-worker/docker'})
