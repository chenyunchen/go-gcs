all:

init: 
	vagrant init bento/centos-7.5

ping:
	export ANSIBLE_CONFIG=../ && ansible all -m ping

ssh:
	ssh -o StrictHostKeyChecking=no vagrant@127.0.0.1 -p 2245 -i ~/.vagrant.d/insecure_private_key