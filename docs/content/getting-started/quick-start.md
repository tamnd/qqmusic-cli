---
title: "Quick start"
description: "Run your first qqmusic command."
weight: 30
---

Once `qqmusic` is on your `PATH`:

```bash
qqmusic --help       # see the command tree
qqmusic version      # build info
```

This is a fresh scaffold, so the command tree is just `version` for now. Add
your first real command in `cli/`, build on the `qqmusic-cli` library package,
and document it here.

A good first command usually fetches one thing and prints it as JSON, so the
output pipes straight into `jq` and the rest of your tools.
