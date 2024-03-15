package main

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/TylerBrock/colorjson"
	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/response"

	"github.com/crossplane/function-generic-resources/input/v1beta1"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

// RunFunction runs the Function.
func (f *Function) RunFunction(_ context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	f.log.Info("Running function", "tag", req.GetMeta().GetTag())

	apiVer := "apiextensions.crossplane.io"
	kind := "environment"
	selector := fmt.Sprintf("%s/%s", apiVer, kind)

	rsp := response.To(req, response.DefaultTTL)

	in := &v1beta1.Input{}
	if err := request.GetInput(req, in); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get Function input from %T", req))
		return rsp, nil
	}

	// TODO: Add your Function logic here!
	response.Normalf(rsp, "I was run with input %q!", in.Example)
	f.log.Info("I was run!", "input", in.Example)

	// // Get function input.
	// in := &v1beta1.Input{}
	// if err := request.GetInput(req, in); err != nil {
	//     response.Fatal(rsp, errors.Wrapf(err, "cannot get Function input from %T", req))
	//     return rsp, nil
	// }
	//

	// // Exit if nothing requested.
	// if in.Spec.EnvironmentConfigs == nil {
	//     f.log.Debug("No EnvironmentConfigs specified, exiting")
	//     return rsp, nil
	// }

	// I think this means get the parent resorucs of the function pipeline, but need to verify.
	// May mean get output XR the pipeline targets.
	oxr, err := request.GetObservedCompositeResource(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get observed composite resource"))
		return rsp, nil
	}

	// // Note(phisco): We need to compute the selectors even if we already
	// // requested them at the previous iteration.
	// requirements, err := buildRequirements(in, oxr)
	// if err != nil {
	//     response.Fatal(rsp, errors.Wrapf(err, "cannot build requirements"))
	//     return rsp, nil
	// }
	//
	requirements, _ := buildRequirements(in, oxr)

	rsp.Requirements = requirements

	// Verify that extra resources were even requested in input.
	if req.ExtraResources == nil {
		f.log.Debug("No extra resources specified, exiting", "requirements", rsp.GetRequirements())
		return rsp, nil
	}

	// Create unstructured object and fetch extra resource into it.
	var inputEnv *unstructured.Unstructured
	if v, ok := request.GetContextKey(req, selector); ok {
		inputEnv = &unstructured.Unstructured{}
		if err := resource.AsObject(v.GetStructValue(), inputEnv); err != nil {
			response.Fatal(rsp, errors.Wrapf(err, "cannot get Composition environment from %T context key %q", req, selector))
			return rsp, nil
		}
		f.log.Debug("Loaded Composition environment from Function context", "context-key", selector)
	}

	extraResources, err := request.GetExtraResources(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get Function input from %T", req))
		return rsp, nil
	}

	fmt.Println(requirements)
	fmt.Println(inputEnv)
	myObj := extraResources["myResource"][0].Resource.Object
	// Make a custom formatter with indent set
	formatter := colorjson.NewFormatter()
	formatter.Indent = 4

	// Marshall the Colorized JSON
	s, _ := formatter.Marshal(myObj)
	fmt.Println(string(s))

	return rsp, nil
}

// Build requirements takes input and outputs an array of external resoruce requirements to request
// from Crossplane's external resource API.
func buildRequirements(_ *v1beta1.Input, _ *resource.Composite) (*fnv1beta1.Requirements, error) {
	extraResources := make(map[string]*fnv1beta1.ResourceSelector, 1) // Define length by input later.
	matchLabels := map[string]string{"type": "cluster"}
	extraResources["myResource"] = &fnv1beta1.ResourceSelector{
		ApiVersion: "apiextensions.crossplane.io/v1alpha1",
		Kind:       "EnvironmentConfig",
		Match: &fnv1beta1.ResourceSelector_MatchLabels{
			MatchLabels: &fnv1beta1.MatchLabels{Labels: matchLabels},
		},
	}
	return &fnv1beta1.Requirements{ExtraResources: extraResources}, nil
}

// // Requirements that must be satisfied for a Function to run successfully.
// type Requirements struct {
//     state         protoimpl.MessageState
//     sizeCache     protoimpl.SizeCache
//     unknownFields protoimpl.UnknownFields
//
//     // Extra resources that this Function requires.
//     // The map key uniquely identifies the group of resources.
//     ExtraResources map[string]*ResourceSelector `protobuf:"bytes,1,rep,name=extra_resources,json=extraResources,proto3" json:"extra_resources,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
// }
// // ResourceSelector selects a group of resources, either by name or by label.
// type ResourceSelector struct {
//     state         protoimpl.MessageState
//     sizeCache     protoimpl.SizeCache
//     unknownFields protoimpl.UnknownFields
//
//     ApiVersion string `protobuf:"bytes,1,opt,name=api_version,json=apiVersion,proto3" json:"api_version,omitempty"`
//     Kind       string `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
//     // Types that are assignable to Match:
//     //
//     //    *ResourceSelector_MatchName
//     //    *ResourceSelector_MatchLabels
//     Match isResourceSelector_Match `protobuf_oneof:"match"`
// }
