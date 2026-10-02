# i3WM Dynamic Tiling Deamon

A lightweight single binary utility written in Go that automatically changes the layout split direction of the **i3 window manager** based on the current window's dimensions.

No dependencies other than go-i3 for ipc interface.

## How to Build

1. Clone or navigate to your project directory.
2. Build the binary using the Go toolchain:
   ```bash
   go build cmd/dynamic-tiling/i3-dynamic-tiling.go
   ```
3. Move the binary to a location in your system path (e.g., `~/.local/bin/` or `/usr/local/bin/` etc).

## How to Use

To have this utility run automatically when i3 starts, add the following line to your i3 configuration file (typically located at `~/.config/i3/config`):

```i3config
# Start the auto-split daemon in the background
exec_always --no-startup-id i3-dynamic-tiling
```

### Alternative: Running with Systemd

If you prefer managing background processes with systemd, you can create a user service instead.
