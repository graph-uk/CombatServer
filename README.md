# Combat server
The server runs tests distributed. It also use "combat" for explore testcases,  "combat-client", for upload tests to server, and "combat-worker" for run tests concurrency, on multiple instances.

# How to develop
The dev process is designed for windows. Linux developers have no any dev-automation yet.

 - Install Git
 - Clone server's repo
 - Place it in `<somePath>\src\github.com\graph-uk\combat-server`.
    It is matters because of make scripts are placing dependencies and tools near
    "src" folder
 - Goto root, run `installCombatDevEnv.cmd` as admin
 
 Now you able:
 
 - Run LiteIDE, and add new amazing features. If you using some other
   IDE - don't forget to set environment vars (`runLiteIDE.cmd`)
 - Build all the apps (combat, client, worker, server), and run
   integration tests (`buildAllAndTestIntegration.cmd`). In progress of testing you can visit http://localhost:9090/sessions/, and watch the progress

If updates have broke some tests - update tests, and re-test.
If updates affect several apps, for example - server, and worker - push updates to both repos.

# How to build/install
On **Linux** you may just run `go get github.com/graph-uk/combat-server && go build github.com/graph-uk/combat-server && go install github.com/graph-uk/combat-server`
On **Windows** - you also need to install MinGW64 for SQLite building.