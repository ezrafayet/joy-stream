# Expose Joy-Stream server with Pinggy (no port forwarding)

Pinggy gives your laptop a public UDP address so clients anywhere can connect without touching your router. **Only the server runs Pinggy;** the client just uses the address Pinggy prints.

## 1. Install Pinggy (on the server machine — your Ubuntu laptop)

1. Download the Linux CLI from **https://pinggy.io/cli/** (choose Linux / amd64 or your arch).
2. Unzip and put the binary somewhere in your path, or in the project folder:
   ```bash
   chmod +x pinggy
   ```
3. (Optional) Sign up at https://dashboard.pinggy.io for a free account if you want a stable URL; otherwise the tunnel still works with a random address.

## 2. Start the Joy-Stream server

In a terminal:

```bash
cd /path/to/joy-stream
./build/server
```

Leave it running. You should see something like: `Joy-Stream UDP server listening on :7355`.

## 3. Start the Pinggy UDP tunnel

In **another** terminal on the **same** machine:

```bash
./pinggy --type udp -l 7355
```

(Use the path to your `pinggy` binary if it’s not in `PATH`.)

Pinggy will print a **public address**, e.g.:

- `udp://a.pinggy.io:12345`, or  
- a hostname and port like `0.tcp.ngrok.io 12345` (example; Pinggy’s format may vary).

Note the **host** and **port** (e.g. `a.pinggy.io` and `12345`).

## 4. Run the client (anywhere)

On the machine that will send gamepad input (same network or across the internet):

```bash
./build/client
```

When prompted for server address, enter the **Pinggy** host and port, e.g.:

- `a.pinggy.io:12345`  
or whatever host and port Pinggy showed (use a colon between host and port).

The client sends UDP to that address; Pinggy forwards it to your laptop’s server. No tunnel or extra software is needed on the client.

## Summary

| Where        | What runs                          |
|-------------|-------------------------------------|
| Server (laptop) | `./build/server` + `./pinggy --type udp -l 7355` |
| Client (any device) | `./build/client` and enter Pinggy host:port     |

If the port changes each time you start Pinggy (common on free tier), use the new host:port in the client. Pro users can reserve a fixed port in the Pinggy dashboard.
