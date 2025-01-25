

go build -ldflags="-s -w -H=windowsgui" #-trimpath
# if -run arg is set, run the program
if [ "$1" = "-run" ]; then
    ./GhostPatcher.exe
fi