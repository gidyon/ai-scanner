rsync:
	rsync -a -e "ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null" ./* sta:/home/ubuntu/sta/scanner

start:
	docker-compose pull && docker-compose up -d --force-recreate

logs:
	docker-compose logs -f --tail=100

clean_dangling:
	docker rmi $(docker images -f dangling=true -q)