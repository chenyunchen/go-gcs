
common-start:
	ansible-playbook -i ../../environments/$(ENV) start.yml

common-stop:
	ansible-playbook -i ../../environments/$(ENV) stop.yml

common-check:
	ansible-playbook -i ../../environments/$(ENV) check.yml

common-ping:
	ansible -i ../../environments/$(ENV) all -m ping

ifeq ($(ENV),)
$(error "Usage: make $(MAKECMDGOALS) ENV=local")
endif
