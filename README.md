Control a virtual Nintendo Switch Joy-Con over the internet to enable local multiplayer on a physical Switch with remote players.

## Building

```bash
go build -o build/server ./server/...
go build -o build/client ./client/...
```

**Linux:** The keyboard client uses gohook and requires X11 dev headers. Install before building:

- Debian/Ubuntu: `sudo apt-get install libx11-dev libx11-xcb-dev libxkbcommon-dev libxkbcommon-x11-dev libxtst-dev`
- Fedora/RHEL: `sudo dnf install libX11-devel`
- Arch: `sudo pacman -S libx11`

(gohook works under X11 only, not Wayland.)