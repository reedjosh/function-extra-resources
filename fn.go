package main

import (
	"context"
	"fmt"

	"github.com/TylerBrock/colorjson"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/function-extra-resources/input/v1beta1"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/response"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

// RunFunction runs the Function.
func (f *Function) RunFunction(_ context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	f.log.Info("Running function", "tag", req.GetMeta().GetTag())

	rsp := response.To(req, response.DefaultTTL)

	// Get function input.
	in := &v1beta1.Input{}
	if err := request.GetInput(req, in); err != nil {
	    response.Fatal(rsp, errors.Wrapf(err, "cannot get Function input from %T", req))
	    return rsp, nil
	}


	// I think this means get the parent resorucs of the function pipeline, but need to verify.
	// May mean get output XR the pipeline targets.
	oxr, err := request.GetObservedCompositeResource(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get observed composite resource"))
		return rsp, nil
	}

	requirements, _ := buildRequirements(in, oxr)

	rsp.Requirements = requirements

	// Verify that extra resources were even requested in input.
	// Appears 100% needed.
	// Still not sure why. This condition should only happen on a bad request.
	if req.ExtraResources == nil {
		f.log.Debug("No extra resources specified, exiting", "requirements", rsp.GetRequirements())
		return rsp, nil
	}

	extraResources, err := request.GetExtraResources(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get Function input from %T", req))
		return rsp, nil
	}

	fmt.Println(requirements)
	myObj := extraResources["myResource"][0].Resource.Object

	// Marshall the Colorized JSON
	// Make a custom formatter with indent set
	formatter := colorjson.NewFormatter()
	formatter.Indent = 4
	s, _ := formatter.Marshal(myObj)
	fmt.Println(string(s))

	return rsp, nil
}

// Build requirements takes input and outputs an array of external resoruce requirements to request
// from Crossplane's external resource API.
func buildRequirements(in *v1beta1.Input, _ *resource.Composite) (*fnv1beta1.Requirements, error) {
	extraResources := make(map[string]*fnv1beta1.ResourceSelector, 1) // Define length by input later.
	matchLabels := map[string]string{"type": "cluster"}
	extraResources["myResource"] = &fnv1beta1.ResourceSelector{
		ApiVersion: in.Spec.ExtraResources[0].APIVersion,
		Kind:       in.Spec.ExtraResources[0].Kind,
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
