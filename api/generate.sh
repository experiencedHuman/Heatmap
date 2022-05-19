protoc AccessPoint.proto \
	--go_out=. \
	--go-grpc_out=. \
	--swift_out=Visibility=Public:../iOS-client/HeatmapUIKit/HeatmapUIKit/ \
	--grpc-swift_out=Visibility=Public,Client=true,Server=false:../iOS-client/HeatmapUIKit/HeatmapUIKit/
