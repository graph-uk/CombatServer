module.exports = class SecEnv {
	constructor()
    {
    	process.env['DOCKER_LOGIN']       = ''
        process.env['DOCKER_PASSWORD']    = ''
    }
};