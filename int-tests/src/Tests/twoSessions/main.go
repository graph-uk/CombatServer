package main

import (
	"Tests_shared/cli"
	"Tests_shared/malibutest"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	theTest := malibutest.NewMalibuTest()
	log.Println(theTest.Params.HostName.Value)

	//Re-create (clear) folders for test binaries
	cli.RemoveAll(`server`)
	cli.RemoveAll(`worker`)
	cli.RemoveAll(`client`)
	cli.MkDir(`server`, 0777)
	cli.MkDir(`worker`, 0777)
	cli.MkDir(`client`, 0777)

	cli.CreateFailTrigger(`testSuccessOrFailure.txt`)

	//Copy compiled binaries to correspond test folders
	cli.CopyFile(`../../Tests_shared/target-app-binaries/malibu-server.exe`, `server/malibu-server.exe`)
	cli.CopyFile(`malibu-server-config.json`, `server/config.json`)
	cli.CopyFile(`../../Tests_shared/target-app-binaries/malibu-worker.exe`, `worker/malibu-worker.exe`)
	cli.CopyFile(`../../Tests_shared/target-app-binaries/malibu-client.exe`, `client/malibu-client.exe`)

	env := cli.EnvRewrite(os.Environ(), `GOPATH`, cli.Pwd()+`/../../Tests_shared/malibuTestsExample`)
	env = cli.EnvRewrite(env, `Path`, cli.Pwd()+`/../../Tests_shared/target-app-binaries`)
	env = cli.EnvExtend(env, `Path`, cli.Pwd()+`/../../../../node_modules/go-win/bin`)

	//env := cli.EnvRewrite(os.Environ(), `GOPATH`, cli.Pwd()+`/../../Tests_shared/malibuTestsExample`)
	//log.Println(env)
	//return

	server := cli.StartCmd(cli.Pwd()+`/server`, &env, `./malibu-server`)
	client := cli.StartCmd(cli.Pwd()+`/../../Tests_shared/malibuTestsExample/src/Tests`, &env, cli.Pwd()+`/client/malibu-client`, `http://localhost:3133`, `./../..`, `40`, `-InternalIP=192.168.1.1`)
	worker := cli.StartCmd(cli.Pwd()+`/worker`, &env, `./malibu-worker`, `http://localhost:3133`)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(`----------------------------------------Server stdout-----------------------------------------`)
			fmt.Println(string(server.StdOutBuf))
			fmt.Println(`----------------------------------------Server stderr-----------------------------------------`)
			fmt.Println(string(server.StdErrBuf))
			fmt.Println(`----------------------------------------Client stdout-----------------------------------------`)
			fmt.Println(string(client.StdOutBuf))
			fmt.Println(`----------------------------------------Client stderr-----------------------------------------`)
			fmt.Println(string(client.StdErrBuf))
			fmt.Println(`----------------------------------------Worker stdout-----------------------------------------`)
			fmt.Println(string(worker.StdOutBuf))
			fmt.Println(`----------------------------------------Worker stderr-----------------------------------------`)
			fmt.Println(string(worker.StdErrBuf))
			panic(r)
		}
	}()

	defer server.Cmd.Process.Kill()
	defer client.Cmd.Process.Kill()
	defer worker.Cmd.Process.Kill()

	// log.Println(env)

	//Check server's output
	//server.WaitingForStdOutContains(`config.json is not found. Default config will be created`, 10*time.Second)
	server.WaitingForStdOutContains(`http server started on`, 10*time.Second)
	server.WaitingForStdOutContains(`Created:  _data/sessions`, 10*time.Second)
	server.WaitingForStdOutContains(`TestFail -InternalIP=192.168.1.1`, 40*time.Second)
	server.WaitingForStdOutContains(`TestSuccess -InternalIP=192.168.1.1`, 40*time.Second)
	server.WaitingForStdOutContains(`Try status: Failed`, 60*time.Second)

	//Check worker's output
	worker.WaitingForStdOutContains(`CaseRunning TestSuccess -InternalIP=192.168.1.1`, 2*time.Minute)
	worker.WaitingForStdOutContains(`Run case... Ok`, time.Minute)
	worker.WaitingForStdOutContains(`CaseRunning TestFail -InternalIP=192.168.1.1`, 2*time.Minute)
	worker.WaitingForStdOutContains(`Run case... Fail`, time.Minute)
	worker.WaitingForStdOutContains(`Failed here for example`, time.Minute)

	//Check client's output
	client.WaitingForStdOutContains(`Cleanup tests`, 10*time.Second)
	client.WaitingForStdOutContains(`Packing tests`, 10*time.Second)
	client.WaitingForStdOutContains(`Uploading session`, 10*time.Second)
	client.WaitingForStdOutContains(` - Pending`, 10*time.Second)
	client.WaitingForStdOutContains(`Case exploring`, time.Minute)
	client.WaitingForStdOutContains(` - Processing`, 40*time.Second)
	client.WaitingForStdOutContains(`Processed 0 of 5 tests`, 40*time.Second)
	client.WaitingForStdOutContains(`Processed 5 of 5 tests`, 400*time.Second)
	client.WaitingForStdOutContains(`Time of testing`, 40*time.Second)
	//return

	// time.Sleep(5 * time.Second)
	// panic(`created`)

	cli.DeleteFailTrigger(`testSuccessOrFailure.txt`)
	client = cli.StartCmd(cli.Pwd()+`/../../Tests_shared/malibuTestsExample/src/Tests`, &env, cli.Pwd()+`/client/malibu-client`, `http://localhost:3133`, `./../..`, `40`, `-InternalIP=192.168.1.1`)
	client.WaitingForStdOutContains(`Time of testing`, 400*time.Second)

	//panic(`test`)

	log.Println(`The test finished succeed.`)

}
