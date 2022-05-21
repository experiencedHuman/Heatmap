protoc -I ./proto \
	--go_out=./proto --go_opt=paths=source_relative \
	--go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=./proto --grpc-gateway_opt=paths=source_relative \
	--swift_out=Visibility=Public:./iOS-client/HeatmapUIKit/HeatmapUIKit/ \
	--grpc-swift_out=Visibility=Public,Client=true,Server=false:./iOS-client/HeatmapUIKit/HeatmapUIKit/ \
	./proto/api/AccessPoint.proto
