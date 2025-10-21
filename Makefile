# Go Mobile Algorand SDK Makefile (Go 1.23.0 compatible)

# ==== Configuration ====
GO := go
PKG := github.com/algorand/go-mobile-algorand-sdk/v2/sdk
OUTPUT_DIR := output

# Resolve gomobile/gobind paths
GOBIN_VAR := $(shell $(GO) env GOBIN)
ifeq ($(strip $(GOBIN_VAR)),)
  BIN := $(shell $(GO) env GOPATH)/bin
else
  BIN := $(GOBIN_VAR)
endif

GOPATH := $(shell $(GO) env GOPATH)
GOMODCACHE := $(shell $(GO) env GOMODCACHE)
GOMOBILE := $(BIN)/gomobile

# Build settings
ANDROID_ABIS := android/arm64,android/amd64
ANDROID_API := 28
IOS_VERSION := 12.0
LD16K := -linkmode=external -extldflags "-Wl,-z,common-page-size=16384 -Wl,-z,max-page-size=16384"

.PHONY: all fmt test clean install-go-mobile android ios verify-android
.DEFAULT_GOAL := all

# ==== Main Targets ====

all: fmt test android ios  ## Format, test, and build both platforms

fmt:  ## Format code
	$(GO) fmt ./...

test:  ## Run tests
	$(GO) test ./... -race

clean:  ## Remove build outputs and caches
	rm -rf $(OUTPUT_DIR)
	rm -rf "$(GOPATH)/pkg/gomobile"
	rm -rf "$(GOMODCACHE)/golang.org/x/mobile"

# ==== Tool Installation ====

install-go-mobile:  ## Install gomobile and gobind (Go 1.23 compatible)
	$(GO) install golang.org/x/mobile/cmd/gomobile@v0.0.0-20230531173138-3c911d8e3eda
	$(GO) install golang.org/x/mobile/cmd/gobind@v0.0.0-20230531173138-3c911d8e3eda
	$(GOMOBILE) init
	$(GOMOBILE) version

# ==== Build Targets ====

android:  ## Build Android AAR
	mkdir -p $(OUTPUT_DIR)
	$(GOMOBILE) bind \
	  -v -trimpath \
	  -target=$(ANDROID_ABIS) \
	  -androidapi $(ANDROID_API) \
	  -o=$(OUTPUT_DIR)/algosdk.aar \
	  -javapkg=com.algorand.algosdk \
	  -ldflags='$(LD16K)' \
	  $(PKG)

ios:  ## Build iOS XCFramework
	mkdir -p $(OUTPUT_DIR)
	$(GOMOBILE) bind \
	  -v -trimpath \
	  -target=ios,iossimulator \
	  -iosversion=$(IOS_VERSION) \
	  -o=$(OUTPUT_DIR)/AlgoSDK.xcframework \
	  -prefix=Algo \
	  $(PKG)

verify-android:  ## Verify Android 16 KB alignment
	@unzip -p $(OUTPUT_DIR)/algosdk.aar jni/arm64-v8a/libgojni.so > /tmp/libgojni.so
	@echo "Checking alignment (expect p_align 0x4000):"
	@llvm-readelf -l /tmp/libgojni.so | grep LOAD || readelf -l /tmp/libgojni.so | grep LOAD