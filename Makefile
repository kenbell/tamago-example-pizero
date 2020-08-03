# Copyright (c) F-Secure Corporation
# https://foundry.f-secure.com
#
# Use of this source code is governed by the license
# that can be found in the LICENSE file.

BUILD_USER = $(shell whoami)
BUILD_HOST = $(shell hostname)
BUILD_DATE = $(shell /bin/date -u "+%Y-%m-%d %H:%M:%S")
BUILD = ${BUILD_USER}@${BUILD_HOST} on ${BUILD_DATE}
REV = $(shell git rev-parse --short HEAD 2> /dev/null)

APP := example-pi-zero
GOENV := GO_EXTLINK_ENABLED=0 CGO_ENABLED=0 GOOS=tamago GOARM=5 GOARCH=arm
TEXT_START := 0x00010000 # Space for interrupt vector, etc

GOFLAGS := -ldflags "-s -w -T $(TEXT_START) -E _rt0_arm_tamago -R 0x1000 -X 'main.Build=${BUILD}' -X 'main.Revision=${REV}'"

SHELL = /bin/bash
JOBS=2

.PHONY: clean

#### primary targets ####

all: $(APP)

elf: $(APP)

install: check_dest elf boot.scr.uimg
	cp boot.scr.uimg $(INSTALLDIR)/boot.scr.uimg
	cp $(APP) $(INSTALLDIR)/$(APP)
	cp config.txt $(INSTALLDIR)/config.txt

#### utilities ####

check_tamago:
	@if [ "${TAMAGO}" == "" ] || [ ! -f "${TAMAGO}" ]; then \
		echo 'You need to set the TAMAGO variable to a compiled version of https://github.com/f-secure-foundry/tamago-go'; \
		exit 1; \
	fi

check_uboot:
	@if [ "${MKIMAGE}" == "" ] || [ ! -f "${MKIMAGE}" ]; then \
		echo 'You need to set the UBOOT variable to a compiled version of mkimage from https://www.denx.de/wiki/U-Boot'; \
		exit 1; \
	fi

check_dest:
	@if [ "${INSTALLDIR}" == "" ] || [ ! -f "${INSTALLDIR}/u-boot-pi0.bin" ]; then \
		echo 'You need to set the INSTALLDIR variable to a mounted Raspberry Pi disk image with u-boot-pi0.bin'; \
		exit 1; \
	fi

clean:
	rm -f $(APP)
	@rm -fr $(APP).bin

#### dependencies ####

$(APP): check_tamago
	$(GOENV) $(TAMAGO) build $(GOFLAGS) -o ${APP}
	$(CROSS_COMPILE)objdump -D $(APP) > $(APP).list

$(APP).bin: $(APP)
	$(CROSS_COMPILE)objcopy -j .text -j .rodata -j .shstrtab -j .typelink \
	    -j .itablink -j .gopclntab -j .go.buildinfo -j .noptrdata -j .data \
	    -j .bss --set-section-flags .bss=alloc,load,contents \
	    -j .noptrbss --set-section-flags .noptrbss=alloc,load,contents\
	    $(APP) -O binary $(APP).bin

boot.scr.uimg: boot.scr check_uboot
	$(MKIMAGE) -A arm -O linux -T script -C none -n boot.scr -d boot.scr boot.scr.uimg
