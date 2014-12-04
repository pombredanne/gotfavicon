all:
	go build
	sudo docker build -t basiclytics/gotfavicon .
	rm gotfavicon
