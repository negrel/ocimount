# OCIMount - Mount OCI/Docker images.

OCIMount is a tool to mount OCI/Docker image easily.

## Getting started

Let's start by installing `ocimount`.

### Installation

```shell
go install github.com/negrel/ocimount@latest
```

### Usage

```shell
# Print help informations
ocimount --help
```

Mount an image as read-only:
```shell
anegrel$ ocimount mount archlinux:latest
INFO[0000] failed to get store, trying again in unshare mode.
ERRO[0000] failed to mount "archlinux:latest": chown /var/home/anegrel/.local/share/containers/storage/overlay/l: operation not permitted

# Oops, it seems that we can't access the storage.
# Let's enter image modified user namespace:
anegrel$ ocimount unshare

root# ocimount mount archlinux:latest
INFO[0000] "docker.io/library/archlinux:latest" successfully mounted at "/var/home/anegrel/.local/share/containers/storage/overlay/de3fc361158be7fbfc230f523b9df392bcf95cba5cf88141292374bf1ec7d2a7/merged".
/var/home/anegrel/.local/share/containers/storage/overlay/de3fc361158be7fbfc230f523b9df392bcf95cba5cf88141292374bf1ec7d2a7/merged

# That's it, out image is mounted read-only.
# Mountpoint is always print to stdout.
```

## Contributing

If you want to contribute to `ocimount` to add a feature or improve the code contact
me at [negrel.dev@protonmail.com](mailto:negrel.dev@protonmail.com), open an
[issue](https://github.com/negrel/ocimount/issues) or make a
[pull request](https://github.com/negrel/ocimount/pulls).

## :stars: Show your support

Please give a :star: if this project helped you!

[![buy me a coffee](.github/images/bmc-button.png)](https://www.buymeacoffee.com/negrel)

## :scroll: License

MIT © [Alexandre Negrel](https://www.negrel.dev/)
