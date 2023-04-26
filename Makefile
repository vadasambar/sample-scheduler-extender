# based on https://stackoverflow.com/a/69400542/6874596
# YYYYMMDD_hhmmss
# $$ is to escape $
# variables in make aren't the same as bash
# for more info: https://stackoverflow.com/a/42462357/6874596
VERSION ?= "latest-$$(date +%Y%m%d_%H%M%S)"

image="ghcr.io/vadasambar/sample-scheduler-extender:${VERSION}"

# Assumes you have logged into GHCR
docker-push: 
	docker build . -t ${image}