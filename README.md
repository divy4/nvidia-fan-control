# nvidia-fan-control

A Golang project to control Nvidia GPU fans.

There aren't any official Linux tools for controlling fan speed.
So here's my take on controlling fans + an excuse to learn more Go.

## Warning

Obviously, anything that controls the speed of your GPU fans can cause your GPU
to overheat if misconfigured or fails unexpectedly.
I've done what I can to ensure whoever uses this doesn't do something really
weird, like slow down the fans when it gets hotter.
But I'm not making any promises you won't break your GPU in the process anyway.
So use at your own risk!

## Build

```bash
make build
```

## Run

```bash
make run
```

## Install

```bash
make install
```

## Uninstall

```bash
make uninstall
```
