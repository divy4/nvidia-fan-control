build:
	cd src && go build -o ../nvidia-fan-control

run: build
	sudo ./nvidia-fan-control run; sudo ./nvidia-fan-control stop
