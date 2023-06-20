# wgnetui

A Fyne-based UI for managing wireguard connections.

## Roadmap

- improve padding in forms
- add GenForm option to purge all devices that don't belong in the generated list (as opposed to just clearing out everything)
  - this can be done by querying all devices that are not contained within the current network
  - might be a little slow though, and has to be done after the generation is done
- when updating the server, offer to regenerate everything over again, if feasible?
- consider adding SSH capabilities, such as:
  - ability to specify SSH hostname (ideally for interoperability with an ssh config file somehow)
- consider adding a ServerDNS field to the genform
- consider setting the peer config to a password multiline entry if it exists
- consider adding a detector/'*' symbol when a device has actually changed
- update the landing screen
- use a splash screen
- future: add a tab for pinging other peers, and doing service checks?
