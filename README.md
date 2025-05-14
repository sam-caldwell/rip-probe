RIP Query Tool
===============

A simple Go program to send a RIP v2 request to a target device (unicast or multicast) and display the routing 
entries received in response. Useful for network diagnostics and discovering RIP-enabled devices.

## Requirements

* Go 1.24+ installed
* Root or elevated privileges (program listens on UDP port 520)

## Installation

```bash
# Clone the repository (or place main.go in your workspace)
git clone <repo-url>
cd <repo-dir>

# Build the binary
go build -o rip-query main.go
```

## Usage

```bash
sudo ./rip-query -target <IP-or-multicast> [-port <udp-port>]
```

| Flag      | Description                                    | Default |
| --------- | ---------------------------------------------- | ------- |
| `-target` | Destination IP or multicast address (required) | —       |
| `-port`   | Destination UDP port for RIP (optional)        | `520`   |

> **Note**: Must be run with sufficient privileges to bind to port 520.

### Examples

* Query a router directly:

  ```bash
  sudo ./rip-query -target 192.168.1.1
  ```

* Query the RIP multicast group:

  ```bash
  sudo ./rip-query -target 224.0.0.9
  ```

* Use a non-standard port:

  ```bash
  sudo ./rip-query -target 224.0.0.9 -port 1520
  ```

## Output

The program prints each RIP route entry with fields:

```
Entry 1: AFI=2, Tag=0, Dest=10.0.0.0, Mask=255.255.255.0, NextHop=0.0.0.0, Metric=1
```

* **AFI**: Address Family Identifier
* **Tag**: Route tag (for OSPF or external routes)
* **Dest**: Destination network
* **Mask**: Subnet mask
* **NextHop**: Next-hop address
* **Metric**: Route metric

## License

MIT License. Feel free to modify and distribute.
