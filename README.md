# Tamago Hello World for Pi

```sh
export TAMAGO=~/work/tamago-go/bin/go
export CROSS_COMPILE=arm-linux-gnueabi-
make example
```

In U-Boot...

## New

boot.scr:

```sh
fatload mmc 0:1 0x10000000 example
bootelf 0x10000000
```

To generate boot.scr.uimg: (`mkimage` is part of U-Boot)

```sh
mkimage -A arm -O linux -T script -C none -n boot.scr -d boot.scr boot.scr.uimg
```

(old)

```sh
fatload mmc 0:1 0x10000 example.bin
go 0x64990
```
