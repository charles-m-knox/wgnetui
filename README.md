# wgnetui

A Fyne-based UI for managing wireguard connections.

## Roadmap

- add GenForm option to purge all devices that don't belong in the generated list (as opposed to just clearing out everything)

- when updating the server, offer to regenerate everything over again, if feasible?
- consider adding a ServerDNS field to the genform
- consider setting the peer config to a password multiline entry if it exists
- add a collapse/show for QR code readout
- consider adding a detector/'*' symbol when a device has actually changed
- implement saving/loading arbitrary sqlite db files
- update the landing screen
- use a splash screen
- future: add a tab for pinging other peers, and doing service checks
