[![Go Report Card](https://goreportcard.com/badge/github.com/darkhz/bluetuith)](https://goreportcard.com/report/github.com/darkhz/bluetuith) [![Packaging status](https://repology.org/badge/tiny-repos/bluetuith.svg)](https://repology.org/project/bluetuith/versions)

![demo](demo/demo.gif)

# bluetuith
bluetuith is a TUI-based bluetooth connection manager, which can interact with bluetooth adapters and devices.
It aims to be a replacement to most bluetooth managers, like blueman.

This is only available on Linux.

This project is currently in the alpha stage.

## Project status
This project has currently been confirmed to be sponsored by the [NLnet](https://nlnet.nl/project/bluetuith/) foundation.
The draft is complete, and the MoU has been signed. The work is now in progress.

Although this repo seems to be currently inactive, please bear in mind that we are actively working on new features, namely:
- Cross-platform support (Windows, MacOS, FreeBSD)
    - Shims[1] for Windows and MacOS
    - Cross platform daemon[2] with a unified API, for any bluetooth app to function across OSes.

- Updating and adding more UI features.
- Extensively refactoring the documentation.

#### Updates
- A new Windows-based shim has been released at [bluetuith-shim-windows](https://github.com/bluetuith-org/bluetuith-shim-windows).

[![Packaging status](https://repology.org/badge/vertical-allrepos/bluetuith.svg)](https://repology.org/project/bluetuith/versions)

## Features
- Transfer and receive files via OBEX.
- Perform pairing with authentication.
- Connect to/disconnect from different devices.
- Interact with bluetooth adapters, toggle power and discovery states
- Connect to or manage Bluetooth based networking/tethering (PANU/DUN)
- Remotely control media playback on the connected device
- Mouse support

## Documentation
The documentation is now hosted [here](https://darkhz.github.io/bluetuith)

The wiki is out-of-date.
