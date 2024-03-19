package validators

import (
	"fmt"
	"sort"

	"github.com/helm-unittest/helm-unittest/pkg/unittest/common"
	"github.com/helm-unittest/helm-unittest/pkg/unittest/valueutils"
)

// LengthEqualDocumentsValidator validate whether the count of manifests rendered form template is Count
type LengthEqualDocumentsValidator struct {
	Paths []string // optional
	Path  string   // optional
	Count int      // optional if paths defined
}

func (v LengthEqualDocumentsValidator) singleValidateCounts(manifest common.K8sManifest, path string, idx, count int) (bool, []string, int) {
	spec, err := valueutils.GetValueOfSetPath(manifest, path)
	if err != nil {
		return false, splitInfof(errorFormat, idx, err.Error()), 0
	}

	if len(spec) == 0 {
		return false, splitInfof(errorFormat, idx, fmt.Sprintf("unknown parameter %s", path)), 0
	}

	specArr, ok := spec[0].([]interface{})
	if !ok {
		return false, splitInfof(errorFormat, idx, fmt.Sprintf("%s is not array", path)), 0
	}
	specLen := len(specArr)
	if count > -1 {
		if specLen != count {
			return false, splitInfof(errorFormat, idx, fmt.Sprintf(
				"count doesn't match as expected. expected: %d actual: %d", count, specLen)), 0
		}
	}
	return true, []string{}, specLen
}

func (v LengthEqualDocumentsValidator) arraysValidateCounts(pathCount map[string]int, idx int) (bool, []string, int) {
	arrayCount := -1

	// Sort alphabetically to get a standardized result
	pathSlice := make([]string, 0)
	for path := range pathCount {
		pathSlice = append(pathSlice, path)
	}

	sort.Strings(pathSlice)

	for _, path := range pathSlice {
		pathCountValue := pathCount[path]
		if arrayCount == -1 {
			arrayCount = pathCountValue
		} else if arrayCount != pathCountValue {
			arrayCount = -1
			return false, splitInfof(errorFormat, idx, fmt.Sprintf(
				"%s count doesn't match as expected. actual: %d", path, pathCountValue)), arrayCount
		}
	}

	return true, []string{}, arrayCount
}

func (v LengthEqualDocumentsValidator) validatePathCount(context *ValidateContext) bool {
	return len(v.Path) > 0 && v.Count == 0
}

func (v LengthEqualDocumentsValidator) validatePathPaths(context *ValidateContext) bool {
	return len(v.Path) > 0 && len(v.Paths) > 0
}

// Validate implement Validatable
func (v LengthEqualDocumentsValidator) Validate(context *ValidateContext) (bool, []string) {
	if v.validatePathCount(context) {
		return false, splitInfof(errorFormat, -1, "'count' field must be set if 'path' is used")
	}
	if v.validatePathPaths(context) {
		return false, splitInfof(errorFormat, -1, "'paths' couldn't be used with 'path'")
	}
	singleMode := len(v.Path) > 0
	manifests, err := context.getManifests()
	if err != nil {
		return false, splitInfof(errorFormat, -1, err.Error())
	}
	validateSuccess := false
	validateErrors := make([]string, 0)
	for idx, manifest := range manifests {
		if singleMode {
			var validateSingleErrors []string
			validateSuccess, validateSingleErrors, _ = v.singleValidateCounts(manifest, v.Path, idx, v.Count)
			validateErrors = append(validateErrors, validateSingleErrors...)
			continue
		} else {
			pathCount := map[string]int{}
			optimizeCheck := true
			for _, path := range v.Paths {
				var validateSingleErrors []string
				validateSuccess, validateSingleErrors, pathCount[path] = v.singleValidateCounts(manifest, path, idx, -1)
				if !validateSuccess {
					validateErrors = append(validateErrors, validateSingleErrors...)
					optimizeCheck = false
				}
			}

			if !optimizeCheck {
				continue
			}

			var arrayCount int
			var validateSingleErrors []string
			validateSuccess, validateSingleErrors, arrayCount = v.arraysValidateCounts(pathCount, idx)
			validateErrors = append(validateErrors, validateSingleErrors...)

			if arrayCount == -1 {
				continue
			}
		}
	}

	if validateSuccess == context.Negative {
		validateSuccess = false
		validateErrors = append(validateErrors, "\texpected result does not match")
	}

	return validateSuccess, validateErrors
}
