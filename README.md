A simple scheduler extender (with `NodeCacheCapable: false`) which filters the nodes that have the label `extender: 'true'` and rejects all other nodes. I have borrowed some code from https://github.com/intel/platform-aware-scheduling/ as is (check the code comments for references).

### Why
I wanted to a simple extender which I could use with scheduler to test https://github.com/kubernetes/autoscaler/pull/5708

### Note
This is extender is only for development and testing purposes. _It is not for production use_. 