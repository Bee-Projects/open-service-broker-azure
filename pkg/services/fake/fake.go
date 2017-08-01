package fake

import (
	"context"

	"github.com/Azure/azure-service-broker/pkg/service"
)

type bindFunction func(
	provisioningContext interface{},
	bindingParameters interface{},
) (interface{}, interface{}, error)

type unbindFunction func(
	provisioningContext interface{},
	bindingContext interface{},
) error

type validationFunction func(bindingParameters interface{}) error

// Module is a fake implementation of the service.Module interface used to
// facilittate testing.
type Module struct {
	ProvisioningValidationBehavior validationFunction
	BindingValidationBehavior      validationFunction
	BindBehavior                   bindFunction
	UnbindBehavior                 unbindFunction
}

// New returns a new instance of a type that fulfills the service.Module
// and provides an example of how such a module is implemented
func New() (*Module, error) {
	return &Module{
		ProvisioningValidationBehavior: defaultValidationBehavior,
		BindingValidationBehavior:      defaultValidationBehavior,
		BindBehavior:                   defaultBindBehavior,
		UnbindBehavior:                 defaultUnbindBehavior,
	}, nil
}

// GetName returns this module's name
func (m *Module) GetName() string {
	return "fake"
}

// ValidateProvisioningParameters validates the provided provisioningParameters
// and returns an error if there is any problem
func (m *Module) ValidateProvisioningParameters(
	provisioningParameters interface{},
) error {
	return m.ProvisioningValidationBehavior(provisioningParameters)
}

// GetProvisioner returns a provisioner that defines the steps a module must
// execute asynchronously to provision a service
func (m *Module) GetProvisioner() (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("run", m.provision),
	)
}

func (m *Module) provision(
	ctx context.Context, // nolint: unparam
	provisioningContext interface{},
	provisioningParameters interface{}, // nolint: unparam
) (interface{}, error) {
	return provisioningContext, nil
}

// ValidateBindingParameters validates the provided bindingParameters and
// returns an error if there is any problem
func (m *Module) ValidateBindingParameters(
	bindingParameters interface{},
) error {
	return m.BindingValidationBehavior(bindingParameters)
}

// Bind synchronously binds to a service
func (m *Module) Bind(
	provisioningContext interface{},
	bindingParameters interface{},
) (interface{}, interface{}, error) {
	return m.BindBehavior(provisioningContext, bindingParameters)
}

// Unbind synchronously unbinds from a service
func (m *Module) Unbind(
	provisioningContext interface{},
	bindingContext interface{},
) error {
	return m.UnbindBehavior(provisioningContext, bindingContext)
}

// GetDeprovisioner returns a deprovisioner that defines the steps a module
// must execute asynchronously to deprovision a service
func (m *Module) GetDeprovisioner() (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("run", m.deprovision),
	)
}

func (m *Module) deprovision(
	ctx context.Context, // nolint: unparam
	provisioningContext interface{},
) (interface{}, error) {
	return provisioningContext, nil
}

func defaultValidationBehavior(params interface{}) error {
	return nil
}

func defaultBindBehavior(
	provisioningContext interface{},
	bindingParameters interface{},
) (interface{}, interface{}, error) {
	return provisioningContext, &Credentials{}, nil
}

func defaultUnbindBehavior(
	provisioningContext interface{},
	bindingContext interface{},
) error {
	return nil
}
