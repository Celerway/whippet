ifeq ($(ARCH),aarch64)
	export GOARCH=arm64
endif
ifeq ($(ARCH),mipsel)
	export GOARCH=mipsle
endif
