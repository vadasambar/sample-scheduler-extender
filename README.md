A simple scheduler extender (with `NodeCacheCapable: false`) which filters the nodes that have the label `extender: 'true'` and rejects all other nodes. I have borrowed some code from https://github.com/intel/platform-aware-scheduling/ as is (check the code comments for references).

### Why
I wanted to a simple extender which I could use with scheduler to test https://github.com/kubernetes/autoscaler/pull/5708

### Note
This is extender is only for development and testing purposes. _It is not for production use_. 

### How to build
If you just want to use the image, you don't need to build it. Use the image from [Packages](https://github.com/vadasambar/sample-scheduler-extender/pkgs/container/sample-scheduler-extender).

To build and push the image,
```
make docker-push REPO=ghcr.io/<username>/sample-scheduler-extender VERSION=v2
```
If you don't specify `VERSION`, Makefile defaults to `latest-<timestamp>` where timestamp is in the format `YYYYMMDD_hhmmss`.
If you don't specify `REPO`, Makefile defaults to `ghcr.io/vadasambar/sample-scheduler-extender`

### Contribution
I look forward to your PRs!