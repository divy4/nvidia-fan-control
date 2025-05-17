build:
	cd src && go build -o ../nvidia-fan-control

run: build
	sudo ./nvidia-fan-control run config.example.json; sudo ./nvidia-fan-control stop config.example.json
