fmt:
	go fmt ./...

test:
	go test ./... -race

install-go-mobile:
	go install golang.org/x/mobile/cmd/gomobile@latest
	go install golang.org/x/mobile/cmd/gobind@latest
	go get golang.org/x/mobile@latest
	gomobile init
	gomobile version

ANDROID_ABIS := android/arm64,android/x86_64
LD16K := -linkmode=external -extldflags "-Wl,--max-page-size=16384 -Wl,-z,max-page-size=16384 -Wl,-z,common-page-size=16384"

android:
	mkdir -p output
	gomobile bind \
	  -target=android/arm64,android/arm,android/386,android/amd64 \
	  -androidapi 28 \
	  -o=output/algosdk.aar \
	  -javapkg=com.algorand.algosdk \
	  -ldflags='-linkmode=external -extldflags "-Wl,-z,max-page-size=16384 -Wl,-z,common-page-size=16384"' \
	  github.com/algorand/go-mobile-algorand-sdk/v2/sdk

ios:
	mkdir -p output
	gomobile bind -target=ios -o=output/AlgoSDK.xcframework -prefix=Algo github.com/algorand/go-mobile-algorand-sdk/v2/sdk

.PHONY: fmt test install-go-mobile android ios