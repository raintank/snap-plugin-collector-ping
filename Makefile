default:
	$(MAKE) deps
	$(MAKE) all
test:
	$(MAKE) deps
	bash -c "./scripts/test.sh $(TEST)"
deps:
	bash -c "./scripts/deps.sh"
check:
	$(MAKE) test
all:
	bash -c "./scripts/build.sh $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))"
