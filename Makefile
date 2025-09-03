build:
	cd src && go build -o ../bin/nvidia-fan-control

run: build
	sudo bin/nvidia-fan-control run config.example.json; sudo bin/nvidia-fan-control stop config.example.json

test:
	cd src && go test

install:
	sudo cp bin/nvidia-fan-control /usr/bin/
	sudo cp systemd/nvidia-fan-control.service /usr/lib/systemd/system/
	sudo systemctl daemon-reload
	echo 'nvidia-fan-control.service installed. Please configure /etc/nvidia-fan-control.json and then enable the service with "systemctl enable --now nvidia-fan-control.service"'

uninstall:
	sudo systemctl disable --now nvidia-fan-control.service
	sudo rm /usr/bin/nvidia-fan-control /usr/lib/systemd/system/nvidia-fan-control.service
	sudo systemctl daemon-reload
