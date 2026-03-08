Control a virtual Nintendo Switch Joy-Con over the internet to enable local multiplayer on a physical Switch with remote players.

## Building

```bash
go build -o build/server ./server/...
go build -o build/client ./client/...
```

**Linux:** The client uses evdev (no X11). Ensure your user can read the keyboard: `sudo adduser $USER input`, then log out and back in.

**Cross-compile for Windows** (from Linux):

```bash
sudo apt-get install gcc-mingw-w64   # one-time
make build-windows
```

Uses MinGW and CGO; the client binary is `build/client.exe`.