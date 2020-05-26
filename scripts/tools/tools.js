const fs = require('fs')

module.exports = class Tools {
	constructor()
    {
        new(require('../env.js'));
		this.mv = require('fs').renameSync
		this.del = require('fs').unlinkSync
		this.cmd = require('child_process').execSync
		this.r = process.env.REPO
    }

  	//env must be stored at ../env.js
	startCmdDetached(cmd, args) {
		var out = require('fs').openSync('.', 'a')
		var child = require('child_process').spawn(cmd, args, { detached: true, stdio: [ 'ignore', out, out ] })
		child.unref()
	}

	checkCmdAvailable(cmd) {
		try {
			this.cmd(cmd)
		}catch{
			console.log('The command "'+cmd+'" returned non-zero code. Check it installed and available')
			console.log('Press any key to exit')
			
			require('fs').readSync(process.stdin.fd, new Buffer(1), 0, 1)
			throw 'check cmd available failed'
		}
	}

	rmdir(dir) {
		if (fs.existsSync(dir)){
			fs.rmdirSync(dir, {recursive: true})
		}
	}

	mkdir(dir) {
		if (!fs.existsSync(dir)){
			fs.mkdirSync(dir)
		}
	}

	mv(oldPath, newPath) {
		fs.renameSync(oldPath,newPath)
	}

	cp(path, newPath) {
		fs.copyFileSync(path,newPath)
	}
}
