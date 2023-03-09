---
description: This section describes the configuration parameters and their types for INX-IRC-Metadata.
keywords:
- IOTA Node 
- Hornet Node
- Metadata
- IRC
- IRC27
- IRC30
- NFT
- Native Token
- Configuration
- JSON
- Customize
- Config
- reference
---


# Core Configuration

INX-IRC-Metadata uses a JSON standard format as a config file. If you are unsure about JSON syntax, you can find more information in the [official JSON specs](https://www.json.org).

You can change the path of the config file by using the `-c` or `--config` argument while executing `inx-irc-metadata` executable.

For example:
```bash
inx-irc-metadata -c config_defaults.json
```

You can always get the most up-to-date description of the config parameters by running:

```bash
inx-irc-metadata -h --full
```

