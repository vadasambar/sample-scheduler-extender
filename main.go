package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"

	"k8s.io/klog/v2"
	extenderv1 "k8s.io/kube-scheduler/extender/v1"
)

// Errors.
var (
	errEmptyBody = errors.New("request body empty")
	errDecode    = errors.New("error decoding request")
)

func filter(w http.ResponseWriter, req *http.Request) {
	extenderArgs := &extenderv1.ExtenderArgs{}

	err := decodeRequest(extenderArgs, req)

	if err != nil {
		klog.Errorf("cannot decode request %v", err)
		w.WriteHeader(http.StatusNotFound)

		return
	}

	filteredNodes := filterNodes(extenderArgs)
	if filteredNodes.Error != "" {
		klog.Error("filtering failed")
		w.WriteHeader(http.StatusNotFound)
	}

	writeResponse(w, filteredNodes)
	klog.Info("filter function done, responded")

}

// filterNodes takes in the arguments for the scheduler and filters nodes based on
// whether the POD resource request fits into each node.
func filterNodes(args *extenderv1.ExtenderArgs) *extenderv1.ExtenderFilterResult {
	var nodeNames []string

	failedNodes := extenderv1.FailedNodesMap{}
	result := extenderv1.ExtenderFilterResult{}

	if args.Nodes == nil || len(args.Nodes.Items) == 0 {
		result.Error = "No nodes to compare. " +
			"This should not happen, perhaps the extender is misconfigured with NodeCacheCapable == true."
		klog.Error(result.Error)

		return &result
	}

	names := []string{}
	for _, node := range args.Nodes.Items {
		names = append(names, node.GetName())
	}

	klog.Infof("filter %v:%v from %v", args.Pod.Namespace, args.Pod.Name, names)

	for _, node := range args.Nodes.Items {
		if node.GetLabels()["extender"] == "true" {
			nodeNames = append(nodeNames, node.GetName())
		} else {

			failedNodes[node.GetName()] = "Doesn't have the label extender='true'"

			continue
		}
	}

	result = extenderv1.ExtenderFilterResult{
		NodeNames:   &nodeNames,
		FailedNodes: failedNodes,
		Error:       "",
	}

	return &result
}

// decodeRequest reads the json request into the given interface args.
// It returns an error if the request is not in the required format.
// based on https://github.com/intel/platform-aware-scheduling/blob/9568b9946134a536a8c9c83f41f42e64fc086b52/gpu-aware-scheduling/pkg/gpuscheduler/scheduler.go#L1308-L1335
func decodeRequest(args interface{}, r *http.Request) error {
	if r.Body == nil {
		return errEmptyBody
	}

	requestDump, dumpErr := httputil.DumpRequest(r, true)
	if dumpErr == nil {
		klog.Infof("http-request:\n%v", string(requestDump))
	}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&args); err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			klog.Errorf("syntax error at byte offset %d", e.Offset)
		}
		return fmt.Errorf("%v: %v", errDecode, err)
	}

	err := r.Body.Close()

	if err != nil {
		err = fmt.Errorf("failed to close request body: %w", err)
	}

	return err
}

// writeResponse takes the incoming interface and writes it as a http response if valid.
// based on https://github.com/intel/platform-aware-scheduling/blob/9568b9946134a536a8c9c83f41f42e64fc086b52/gpu-aware-scheduling/pkg/gpuscheduler/scheduler.go#L1337-L1343
func writeResponse(w http.ResponseWriter, result interface{}) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(result); err != nil {
		http.Error(w, "Encode error", http.StatusBadRequest)
	}
}

func main() {

	http.HandleFunc("/apiv1/filter", filter)

	klog.Info("Starting server at port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
