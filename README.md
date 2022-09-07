# CHANGELOG.md Management Tool
## *cl-tool*

Manage your `CHANGELOG.md` in your team without merge conflicts!

This is a helper tool that allows your changelog data to be updated in concurrent Pull Requests without having to fix merge conflicts for each merged PR.

## How it works

Your `CHANGELOG.md` file is generated from entries stored in a `.changelog` folder.
The folder structure of the `.changelog` folder is `.changelog/<version>/<entry>.yml`.
Each yaml file under a version can contain one or multiple entries.

```yaml
# 2022-05-30-hello-world.yml

added:
  - 'this is a new feature'
changed:
  - 'changed feature 1'
  - 'changed feature 2'
```

Either add entries under a known version, or add all features under the `.changelog/Unreleased` folder.

At release time, run the `release` command to rename the `Unreleased` folder with the new version and regenerate the changelog with the `generate` command.

![](/images/cl-tool.gif)

## Commands

- `new`: Create new changelog entries, under the `.changelog/Unreleased` folder
- `generate`: Output a generated CHANGELOG
- `release`: Move `.changelog/Unreleased` entries into a new versioned/dated folder.

Inspired by:
- [github.com/nettsundere/cyberlog](https://github.com/nettsundere/cyberlog)
- [github.com/uptech/git-cl](https://github.com/uptech/git-cl)

## Custom templates

Add a `.changelog/head.md` file to override the default header text. Add a `.changelog/tail.md` file to add text at the bottom of the changelog.
