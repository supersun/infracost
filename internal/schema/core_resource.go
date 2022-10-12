package schema

type CoreResourceFunc func(*ResourceData) CoreResource

// CoreResource is the new/preferred way to represent provider-agnostic resources that
// support advanced features such as Infracost Cloud usage estimates and actual costs.
type CoreResource interface {
	CoreType() string
	UsageSchema() []*UsageItem
	PopulateUsage(u *UsageData)
	BuildResource() *Resource
}

type UsageParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CoreResourceWithUsageParams is a CoreResource that sends additional
// parameters (e.g. Lambda memory size) to the Usage API when estimating usage.
type CoreResourceWithUsageParams interface {
	CoreResource
	UsageEstimationParams() []UsageParam
}

// PartialResource is used to collect all information required to construct a
// resource that is generated by provider parser and pass it back up to
// top level functions that can supply additional provider-agnostic information
// (such as Infracost Cloud usage estimates) before the resource is built.
type PartialResource struct {
	ResourceData *ResourceData

	// CoreResource is the new/preferred struct for providing an intermediate-object
	// that contains all provider-derived information, but has not yet been built into
	// a Resource.
	CoreResource CoreResource

	// Resource field is provided for backward compatibility with provider resource builders
	// that have not yet been converted to build CoreResource's
	Resource *Resource

	// CloudResourceIDs are collected during parsing in case they need to be uploaded to the
	// Cloud Usage API to be used in the usage estimate calculations.
	CloudResourceIDs []string
}

// BuildResource create a new Resource from the CoreResource, or (for backward compatibility) returns
// a previously built Resource
func BuildResource(partial *PartialResource, fetchedUsage *UsageData) *Resource {
	var res *Resource
	if partial.CoreResource != nil {
		u := partial.ResourceData.UsageData
		u = u.Merge(fetchedUsage)

		partial.CoreResource.PopulateUsage(u)
		res = partial.CoreResource.BuildResource()
	} else {
		res = partial.Resource
	}

	if res != nil {
		res.ResourceType = partial.ResourceData.Type
		res.Tags = partial.ResourceData.Tags
		res.Metadata = partial.ResourceData.Metadata
	}

	return res
}
