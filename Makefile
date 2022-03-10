APP_NAME=heatmap
BINARY_NAME=${APP_NAME}.out

build:
	go build -o ${BINARY_NAME} ${APP_NAME}.go

run:
	go build -o ${BINARY_NAME} ${APP_NAME}.go
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}
