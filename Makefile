include $(TOPDIR)/rules.mk

PKG_NAME:=whippet
PKG_VERSION:=2024-06-07
PKG_RELEASE=0.1

PKG_BUILD_DEPENDS:=golang

#GO_PKG:=github.com/celerway/whippet
#GO_PKG_BUILD_PKG:=github.com/celerway/whippet

CMAKE_INSTALL:=0
PKG_USE_NINJA:=0

PKG_LICENSE:=BSD-3
PKG_LICENSE_FILES:=LICENSE.md

PKG_BUILD_PARALLEL:=1

include $(INCLUDE_DIR)/package.mk
include $(TOPDIR)/feeds/packages/lang/golang/golang-package.mk

define Build/Prepare
	mkdir -p $(PKG_BUILD_DIR)
	cp -r src/* $(PKG_BUILD_DIR)/
endef

define Package/whippet
  SECTION:=celerway
  CATEGORY:=Celerway
  DEPENDS:=$(GO_ARCH_DEPENDS)
  TITLE:=MQTTv5 request/response command-line tool
endef

TARGET_CFLAGS += -I$(STAGING_DIR)/usr/include

define Package/whippet/install
	$(INSTALL_DIR) $(1)/usr/bin
	$(CP) $(PKG_BUILD_DIR)/whippet-cli $(1)/usr/bin/whippet
endef

$(eval $(call GoBinPackage,whippet))
$(eval $(call BuildPackage,whippet))
