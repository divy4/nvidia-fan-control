build:
	mkdir build/
	cd src && go build -o ../build/nvidia-fan-control

run: build
	sudo build/nvidia-fan-control run config.example.json; sudo build/nvidia-fan-control stop config.example.json

package: build
