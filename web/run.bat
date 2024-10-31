go build -o cloudpilot-build.exe cmd/web/.
cloudpilot-build.exe -dbname=cloudpilot -dbuser=postgres -cache=false -production=false