package main

import (
	"context"
	"fmt"

	"github.com/TylerBrock/colorjson"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/response"

	"github.com/crossplane/function-extra-resources/input/v1beta1"
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

	// // Get XR the pipeline targets.
	// oxr, err := request.GetObservedCompositeResource(req)
	// if err != nil {
	// 	response.Fatal(rsp, errors.Wrapf(err, "cannot get observed composite resource"))
	// 	return rsp, nil
	// }


	requirements, _ := buildRequirements(in)
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
		response.Fatal(rsp, errors.Wrapf(err, "fetching extra resources %T", req))
		return rsp, nil
	}

	// Remove eventually. Just here to output learnings for now.
	for _, objs := range(extraResources) {
		for _, obj := range(objs) {
			pretty(obj.Resource.Object)
		}
	}

	return rsp, nil
}

func pretty(obj interface{}) {
	// Marshall the Colorized JSON
	// Make a custom formatter with indent set
	formatter := colorjson.NewFormatter()
	formatter.Indent = 4
	s, _ := formatter.Marshal(obj)
	fmt.Println(string(s))
}

// Build requirements takes input and outputs an array of external resoruce requirements to request
// from Crossplane's external resource API.
func buildRequirements(in *v1beta1.Input) (*fnv1beta1.Requirements, error) {
	extraResources := make(map[string]*fnv1beta1.ResourceSelector, len(in.Spec.ExtraResources)) // Define length by input later.
	for _, resourceRequest := range in.Spec.ExtraResources {

		matchLabels := map[string]string{"type": "cluster"}
		extraResources[resourceRequest.Into] = &fnv1beta1.ResourceSelector{
			ApiVersion: resourceRequest.APIVersion,
			Kind:       resourceRequest.Kind,
			Match: &fnv1beta1.ResourceSelector_MatchLabels{
				MatchLabels: &fnv1beta1.MatchLabels{Labels: matchLabels},
			},
		}
	}
	return &fnv1beta1.Requirements{ExtraResources: extraResources}, nil
}
