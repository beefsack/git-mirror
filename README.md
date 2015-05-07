# git-mirror - simple Git mirrors

`git-mirror` is designed to create and serve read-only mirrors of your Git repositories locally or wherever you choose.  A recent GitHub outage reinforces the fact that developers shouldn't be relying on a single remote for hosting code.

A major design goal of `git-mirror` is that it should just work with as little configuration as possible.

## Get started

Download and extract the latest release from the [releases page](https://github.com/beefsack/git-mirror/releases).

Create `config.toml` similar to:

```toml
[[repo]]
Origin = "https://github.com/beefsack/git-mirror.git"
```
By default it will update the mirror every **15 minutes** and will serve the mirror over HTTP using port **8080**.  You can specify as many repos as you want by having multiple `[[repo]]` sections.

Run `git-mirror` with the path to the config file:

```bash
$ ./git-mirror config.toml
2015/05/07 11:08:06 starting web server on :8080
2015/05/07 11:08:06 updating github.com/beefsack/git-mirror.git
2015/05/07 11:08:08 updated github.com/beefsack/git-mirror.git
```

Now you can clone from your mirror on the default port of `8080`:

```bash
$ git clone http://localhost:8080/github.com/beefsack/git-mirror.git
Cloning into 'git-mirror'...
Checking connectivity... done.
```

## Advanced configuration

See [the example config](example-config.toml) for more advanced configurations.

## Authentication and authorisation

If you wish to control access to the mirror or specific repositories, consider proxying to `git-mirror` using a web server such as Nginx.
