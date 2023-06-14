# wgnetui

A Fyne-based UI for managing wireguard connections.

## Requirements

Builds upon the requirements for `wgnetgen`.

- all flags from `wgnetgen` must be translated to form input fields in this

- add a "description" and "in use" flags to each config
- store everything into a sqlite gorm db
- add & track clients
- modified config in addition to auto-generated config
  - later on: versioned changes to auto-generated configs?
- QR code reader for generated configs
  - generate png on the fly and render with Fyne
