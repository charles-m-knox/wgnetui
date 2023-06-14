# About

Wireguard Network UI (`wgnetui`) can generate an entire network of connected devices with a layer of customizability added on top.

It works by taking an IPv4 address range such as `192.168.5.0/24` and generating a Wireguard peer config for every single IPv4 address within that range.

A single server acts as the "hub" for accepting connections from all other peers on the mesh network. This server exists within the range `192.168.5.0/24` and has a large Wireguard config file that contains a peer connection for each peer.

## Behavior

By default, generated Wireguard public and private keys are kept between configuration changes. Keys are assigned in a sequential order to whichever devices you choose first.

As a simplified example, let's assume your desired network will only have 5 devices currently, and you are operating over the IP address range `192.168.5.0/24`, leaving plenty of room for more devices in the future.

```text
Device 1: main-server (primary wireguard server)
Device 2: my-laptop
Device 3: my-phone
Device 4: my-data-store
Device 5: my-raspberry-pi

... the remaining IP addresses are unallocated for now, and that's OK
```

With the above device setup, `wgnetui` will assign the first 5 IP addresses of *any* IP address CIDR range to these devices - in this case, the allocation will occur as follows:

```text
192.168.5.1: main-server (primary wireguard server)
192.168.5.2: my-laptop
192.168.5.3: my-phone
192.168.5.4: my-data-store
192.168.5.5: my-raspberry-pi
```

Now, let's suppose you actually wanted to change your CIDR range to something different, such as `192.168.25.0/23`, because you've realized you want a bigger subnet and a different IP range.

This is where the key persistence feature comes into play - it makes it just a bit easier to change your existing Wireguard clients that are already configured without messing with new keys. Keys will stay familiar to you over time instead of changing constantly.

This behavior can easily be overridden to generate new keys on your next config update by toggling "Regenerate Keys".
