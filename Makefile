init-data:
	./scripts/startup_container.sh ORABFILE 12345
	cp ./scripts/load_data.sh ./tmp/load_data.sh
	docker exec -it mentor-ora /bin/bash -c "/tmp/load_data.sh"

clean:
	docker rm -f mentor-ora
	rm -rf ./tmp
