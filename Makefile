
fmt:
	go fmt ./...

test:
	go test ./... -race

install-go-mobile:
	go install golang.org/x/mobile/cmd/gomobile@v0.0.0-20240909163608-642950227fb3
	go get -d golang.org/x/mobile@v0.0.0-20240909163608-642950227fb3

	gomobile init
	gomobile version

android:
	mkdir -p output
	gomobile bind -target=android -o=output/algosdk.aar -javapkg=com.algorand.algosdk github.com/algorand/go-mobile-algorand-sdk/v2/sdk

ios:
	mkdir -p output
	gomobile bind -target=ios -o=output/AlgoSDK.xcframework -prefix=Algo github.com/algorand/go-mobile-algorand-sdk/v2/sdk

.PHONY: fmt test install-go-mobile android ios
