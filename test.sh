go build -ldflags="-s -w" -o GhostLauncher.exe
rm -f xHydra/GhostLauncher.exe
cp GhostLauncher.exe xHydra/GhostLauncher.exe
cd xHydra
./GhostLauncher.exe
cd ..