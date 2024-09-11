PWD := $(shell pwd)
export GO1111MODULE=on

config ?= example.yaml
case ?= test

args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

api:
	go install
	@apiman --config=$(config) --work=$(PWD)/example run

run:
	go install
	@apiman --config=$(config) --verbose --work=$(PWD)/example case $(case)
	# @apiman --config=$(config) --work=$(PWD)/example case $(case)

bench:
	go install
	@apiman --config=$(config) --work=$(PWD)/example case --bench=100 --num-worker=10 --interval=1 $(case)