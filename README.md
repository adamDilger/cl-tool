# Changelog Tool

Manage your `CHANGELOG.md` in your team without merge conflicts!

This is a helper tool that allows your changelog data to be updated in concurrent Pull Requests without having to fix merge conflicts.

## How it works

Changelog entries are stored as `yaml` files in a `.changelog` folder. These files can contain one or multiple changelog entries and are stored in the folder structure of `.changelog/<version>/<entry>`

```yaml
# 2022-05-30-hello-world.yml

added:
  - 'this is a new feature'
changed:
  - 'changed feature 1'
  - 'changed feature 2'
```

Either add entries under a known version, or add all features under the `.changelog/Unreleased` folder.

At release time, rename the `Unreleased` folder with the appropriate version number + date and regenerate the changelog.

![](/images/generate.jpg)

## Commands

- `new`: Create new changelog entries, under the `.changelog/Unreleased` folder
- `generate`: Output a generated CHANGELOG
- `release`: Performs multiple commands that:
  - Requests the verison number from the user e.g. 2.1.3
  - Renames `.changelog/Unreleased` folder with the current date and version
  - Regenerates the `CHANGELOG.md`
  - Updates the version number in the `pom.xml`
    - TODO: this should be an opt in feature
  - Creates a git branch and commit with changes

Inspired by:
  - [github.com/nettsundere/cyberlog](https://github.com/nettsundere/cyberlog)
  - [github.com/uptech/git-cl](https://github.com/uptech/git-cl)

## Custom templates

Add a `.changelog/head.md` file to override the default header text. Add a `.changelog/tail.md` file to add text at the bottom of the changelog.

## Creating a release

1. Update the .changelog
2. `make compile version=<new-verison>`
2. create `<new-version>` tag