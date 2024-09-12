PWD := $(shell pwd)
export GO1111MODULE=on

config ?= example.yaml
case ?= test

args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

api:
	go install
	@apiman --config=$(config) --work=$(PWD)/example api

run:
	go install
	# @apiman --config=$(config) --verbose --work=$(PWD)/example run $(case)
	@apiman --config=$(config) --work=$(PWD)/example run $(case)

bench:
	go install
	@apiman --config=$(config) --work=$(PWD)/example bench $(case) --num-bench=50 --num-worker=10 --interval=0