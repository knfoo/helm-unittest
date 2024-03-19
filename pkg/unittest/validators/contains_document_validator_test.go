package validators_test

import (
	"testing"

	"github.com/helm-unittest/helm-unittest/pkg/unittest/common"
	. "github.com/helm-unittest/helm-unittest/pkg/unittest/validators"
	"github.com/stretchr/testify/assert"
)

var docToTestContainsDocument1 = `
apiVersion: v1
kind: Service
metadata:
  name: foo
  namespace: bar
`

var docToTestContainsDocument2 = `
apiVersion: v1
kind: Service
metadata:
  name: bar
  namespace: foo
`

var docToTestContainsDocument3 = `
apiVersion: v1
kind: Service
metadata:
    name: bar	
`

var docToTestContainsDocument4 = `
apiVersion: v1
kind: Service
metadata:
    namespace: foo
`

func TestContainsDocumentValidatorWhenEmptyNOk(t *testing.T) {
	validator := ContainsDocumentValidator{
		Kind:       "Service",
		APIVersion: "v1",
		Name:       "bar",
		Namespace:  "foo",
		Any:        true,
	}
	pass, diff := validator.Validate(&ValidateContext{
		Index: -1,
		Docs:  []common.K8sManifest{},
	})

	assert.False(t, pass)
	assert.Equal(t, []string{
		"DocumentIndex:	0",
		"Expected to contain document:",
		"\tKind = Service, apiVersion = v1, Name = bar, Namespace = foo",
	}, diff)
}

func TestContainsDocumentValidatorNegativeWhenEmptyOk(t *testing.T) {
	validator := ContainsDocumentValidator{
		Kind:       "Service",
		APIVersion: "v1",
		Name:       "bar",
		Namespace:  "foo",
		Any:        true,
	}
	pass, diff := validator.Validate(&ValidateContext{
		Index:    -1,
		Docs:     []common.K8sManifest{},
		Negative: true,
	})

	assert.False(t, pass)
	assert.Equal(t, []string{}, diff)
}

func TestContainsDocumentValidatorWhenNotAllDocumentsAreOk(t *testing.T) {
	validator := ContainsDocumentValidator{
		Kind:       "Service",
		APIVersion: "v1",
		Name:       "bar",
		Namespace:  "foo",
	}
	pass, diff := validator.Validate(&ValidateContext{
		Index: -1,
		Docs: []common.K8sManifest{makeManifest(docToTestContainsDocument1),
			makeManifest(docToTestContainsDocument2)},
	})

	assert.False(t, pass)
	assert.Equal(t, []string{
		"DocumentIndex:\t0",
		"Expected to contain document:",
		"\tKind = Service, apiVersion = v1, Name = bar, Namespace = foo",
	}, diff)
}

func TestContainsDocumentValidatorWhenAtleastOneDocumentsIsOk(t *testing.T) {
	validator := ContainsDocumentValidator{
		Kind:       "Service",
		APIVersion: "v1",
		Name:       "foo",
		Namespace:  "bar",
		Any:        true,
	}
	pass, diff := validator.Validate(&ValidateContext{
		Index: -1,
		Docs: []common.K8sManifest{makeManifest(docToTestContainsDocument1),
			makeManifest(docToTestContainsDocument2)},
	})

	assert.True(t, pass)
	assert.Equal(t, []string{}, diff)
}

func TestContainsDocumentValidatorWhenAtleastOneDocumentsIsOkInverse(t *testing.T) {
	validator := ContainsDocumentValidator{
		Kind:       "Service",
		APIVersion: "v1",
		Name:       "bar",
		Namespace:  "foo",
		Any:        true,
	}
	pass, diff := validator.Validate(&ValidateContext{
		Index: -1,
		Docs: []common.K8sManifest{makeManifest(docToTestContainsDocument1),
			makeManifest(docToTestContainsDocument2)},
		Negative: true,
	})

	assert.False(t, pass)
	assert.Equal(t, []string{
		"DocumentIndex:\t1",
		"Expected NOT to contain document:",
		"\tKind = Service, apiVersion = v1, Name = bar, Namespace = foo",
	}, diff)
}

func TestContainsDocumentValidatorIndexWhenOk(t *testing.T) {
	validator := ContainsDocumentValidator{
		Kind:       "Service",
		APIVersion: "v1",
		Name:       "bar",
		Namespace:  "foo",
		Any:        false,
	}
	pass, diff := validator.Validate(&ValidateContext{
		Index: 1,
		Docs: []common.K8sManifest{makeManifest(docToTestContainsDocument1),
			makeManifest(docToTestContainsDocument2)},
	})

	assert.True(t, pass)
	assert.Equal(t, []string{}, diff)
}

func TestContainsDocumentValidatorNoNameWhenOk(t *testing.T) {
	validator := ContainsDocumentValidator{
		Kind:       "Service",
		APIVersion: "v1",
		Namespace:  "foo",
		Any:        false,
	}

	pass, diff := validator.Validate(&ValidateContext{
		Index: -1,
		Docs:  []common.K8sManifest{makeManifest(docToTestContainsDocument2)},
	})

	assert.True(t, pass)
	assert.Equal(t, []string{}, diff)
}

func TestContainsDocumentValidatorNoNamespaceWhenOk(t *testing.T) {
	validator := ContainsDocumentValidator{
		Kind:       "Service",
		APIVersion: "v1",
		Name:       "foo",
		Namespace:  "",
		Any:        false,
	}

	pass, diff := validator.Validate(&ValidateContext{
		Index: -1,
		Docs:  []common.K8sManifest{makeManifest(docToTestContainsDocument1)},
	})

	assert.True(t, pass)
	assert.Equal(t, []string{}, diff)
}

func TestContainsDocumentValidatorNoNamespaceWhenNegativeOk(t *testing.T) {
	validator := ContainsDocumentValidator{
		Kind:       "InvalidService",
		APIVersion: "v1",
		Name:       "foo",
		Namespace:  "",
		Any:        true,
	}

	pass, diff := validator.Validate(&ValidateContext{
		Index:    -1,
		Docs:     []common.K8sManifest{makeManifest(docToTestContainsDocument1)},
		Negative: true,
	})

	assert.True(t, pass)
	assert.Equal(t, []string{}, diff)
}

func TestContainsDocumentValidatorNoNameNamespaceWhenOk(t *testing.T) {
	validator := ContainsDocumentValidator{
		Kind:       "Service",
		APIVersion: "v1",
		Name:       "",
		Namespace:  "",
		Any:        false,
	}

	pass, diff := validator.Validate(&ValidateContext{
		Index: -1,
		Docs: []common.K8sManifest{makeManifest(docToTestContainsDocument1),
			makeManifest(docToTestContainsDocument2)},
	})

	assert.True(t, pass)
	assert.Equal(t, []string{}, diff)
}

func TestContainsDocumentValidatorNoNameNamespaceWhenNegativeNOk(t *testing.T) {
	validator := ContainsDocumentValidator{
		Kind:       "Service",
		APIVersion: "v1",
		Name:       "",
		Namespace:  "",
		Any:        false,
	}

	pass, diff := validator.Validate(&ValidateContext{
		Index: -1,
		Docs: []common.K8sManifest{makeManifest(docToTestContainsDocument1),
			makeManifest(docToTestContainsDocument2)},
		Negative: true,
	})

	assert.False(t, pass)
	assert.Equal(t, []string{
		"DocumentIndex:\t0",
		"Expected NOT to contain document:",
		"\tKind = Service, apiVersion = v1, Name = , Namespace =",
		"DocumentIndex:\t1",
		"Expected NOT to contain document:",
		"\tKind = Service, apiVersion = v1, Name = , Namespace =",
	}, diff)
}

func TestContainsDocumentValidatorWhenFailKind(t *testing.T) {
	validator := ContainsDocumentValidator{
		Kind:       "Deployment",
		APIVersion: "apps/v1",
		Name:       "foo",
		Namespace:  "bar",
		Any:        false,
	}

	pass, diff := validator.Validate(&ValidateContext{
		Index: -1,
		Docs: []common.K8sManifest{makeManifest(docToTestContainsDocument1),
			makeManifest(docToTestContainsDocument2)},
	})

	assert.False(t, pass)
	assert.Equal(t, []string{
		"DocumentIndex:\t0",
		"Expected to contain document:",
		"\tKind = Deployment, apiVersion = apps/v1, Name = foo, Namespace = bar",
		"DocumentIndex:\t1",
		"Expected to contain document:",
		"\tKind = Deployment, apiVersion = apps/v1, Name = foo, Namespace = bar",
	}, diff)
}

func TestContainsDocumentValidatorWhenFailAPIVersion(t *testing.T) {
	validator := ContainsDocumentValidator{
		Kind:       "Service",
		APIVersion: "apps/v1",
		Name:       "foo",
		Namespace:  "bar",
		Any:        false,
	}

	pass, diff := validator.Validate(&ValidateContext{
		Index: -1,
		Docs: []common.K8sManifest{makeManifest(docToTestContainsDocument1),
			makeManifest(docToTestContainsDocument2)},
	})

	assert.False(t, pass)
	assert.Equal(t, []string{
		"DocumentIndex:\t0",
		"Expected to contain document:",
		"\tKind = Service, apiVersion = apps/v1, Name = foo, Namespace = bar",
		"DocumentIndex:\t1",
		"Expected to contain document:",
		"\tKind = Service, apiVersion = apps/v1, Name = foo, Namespace = bar",
	}, diff)
}

func TestContainsDocumentValidatorFail(t *testing.T) {
	tests := []struct {
		name           string
		validator      ContainsDocumentValidator
		fixtureContext ValidateContext
		expected       []string
	}{
		{
			name: "it should not fail when namespace is not specified....",
			validator: ContainsDocumentValidator{
				Kind:       "Service",
				APIVersion: "apps/v1",
				Name:       "foo",
			},
			fixtureContext: ValidateContext{
				Index: 0,
				Docs:  []common.K8sManifest{makeManifest(docToTestContainsDocument3)},
			},
			expected: []string{
				"DocumentIndex:\t0",
				"Expected to contain document:",
				"\tKind = Service, apiVersion = apps/v1, Name = foo, Namespace =",
			},
		},
		{
			name: "it should not fail when name is not specified....",
			validator: ContainsDocumentValidator{
				Kind:       "Service",
				APIVersion: "apps/v1",
				Namespace:  "bar",
			},
			fixtureContext: ValidateContext{
				Index: 0,
				Docs:  []common.K8sManifest{makeManifest(docToTestContainsDocument4)},
			},
			expected: []string{
				"DocumentIndex:\t0",
				"Expected to contain document:",
				"\tKind = Service, apiVersion = apps/v1, Name = , Namespace = bar",
			},
		},
	}

	for _, test := range tests {
		pass, diff := test.validator.Validate(&test.fixtureContext)
		assert.False(t, pass)
		assert.Equal(t, test.expected, diff)
	}
}
