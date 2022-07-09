package v1alpha1

import (
	"encoding/json"

	"github.com/ghodss/yaml"
)

// KnownTypeField contains mapping between CRD field and known Kubernetes type.
// This is mainly used for unit conversion in unknown resources (e.g. 0.1 == 100mi)
// TODO: Describe the members of this type
type KnownTypeField struct {
	Field string `json:"field,omitempty" protobuf:"bytes,1,opt,name=field"`
	Type  string `json:"type,omitempty" protobuf:"bytes,2,opt,name=type"`
}

type OverrideIgnoreDiff struct {
	//JSONPointers is a JSON path list following the format defined in RFC4627 (https://datatracker.ietf.org/doc/html/rfc6902#section-3)
	JSONPointers []string `json:"jsonPointers" protobuf:"bytes,1,rep,name=jSONPointers"`
	//JQPathExpressions is a JQ path list that will be evaludated during the diff process
	JQPathExpressions []string `json:"jqPathExpressions" protobuf:"bytes,2,opt,name=jqPathExpressions"`
	// ManagedFieldsManagers is a list of trusted managers. Fields mutated by those managers will take precedence over the
	// desired state defined in the SCM and won't be displayed in diffs
	ManagedFieldsManagers []string `json:"managedFieldsManagers" protobuf:"bytes,3,opt,name=managedFieldsManagers"`
}

type rawResourceOverride struct {
	HealthLua         string           `json:"health.lua,omitempty"`
	UseOpenLibs       bool             `json:"health.lua.useOpenLibs,omitempty"`
	Actions           string           `json:"actions,omitempty"`
	IgnoreDifferences string           `json:"ignoreDifferences,omitempty"`
	KnownTypeFields   []KnownTypeField `json:"knownTypeFields,omitempty"`
}

// ResourceOverride holds configuration to customize resource diffing and health assessment
// TODO: describe the members of this type
type ResourceOverride struct {
	HealthLua         string             `protobuf:"bytes,1,opt,name=healthLua"`
	UseOpenLibs       bool               `protobuf:"bytes,5,opt,name=useOpenLibs"`
	Actions           string             `protobuf:"bytes,3,opt,name=actions"`
	IgnoreDifferences OverrideIgnoreDiff `protobuf:"bytes,2,opt,name=ignoreDifferences"`
	KnownTypeFields   []KnownTypeField   `protobuf:"bytes,4,opt,name=knownTypeFields"`
}

// TODO: describe this method
func (s *ResourceOverride) UnmarshalJSON(data []byte) error {
	raw := &rawResourceOverride{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	s.KnownTypeFields = raw.KnownTypeFields
	s.HealthLua = raw.HealthLua
	s.UseOpenLibs = raw.UseOpenLibs
	s.Actions = raw.Actions
	return yaml.Unmarshal([]byte(raw.IgnoreDifferences), &s.IgnoreDifferences)
}

// TODO: describe this method
func (s ResourceOverride) MarshalJSON() ([]byte, error) {
	ignoreDifferencesData, err := yaml.Marshal(s.IgnoreDifferences)
	if err != nil {
		return nil, err
	}
	raw := &rawResourceOverride{s.HealthLua, s.UseOpenLibs, s.Actions, string(ignoreDifferencesData), s.KnownTypeFields}
	return json.Marshal(raw)
}

// TODO: describe this method
func (o *ResourceOverride) GetActions() (ResourceActions, error) {
	var actions ResourceActions
	err := yaml.Unmarshal([]byte(o.Actions), &actions)
	if err != nil {
		return actions, err
	}
	return actions, nil
}

// TODO: describe this type
// TODO: describe members of this type
type ResourceActions struct {
	ActionDiscoveryLua string                     `json:"discovery.lua,omitempty" yaml:"discovery.lua,omitempty" protobuf:"bytes,1,opt,name=actionDiscoveryLua"`
	Definitions        []ResourceActionDefinition `json:"definitions,omitempty" protobuf:"bytes,2,rep,name=definitions"`
}

// TODO: describe this type
// TODO: describe members of this type
type ResourceActionDefinition struct {
	Name      string `json:"name" protobuf:"bytes,1,opt,name=name"`
	ActionLua string `json:"action.lua" yaml:"action.lua" protobuf:"bytes,2,opt,name=actionLua"`
}

// TODO: describe this type
// TODO: describe members of this type
type ResourceAction struct {
	Name     string                `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Params   []ResourceActionParam `json:"params,omitempty" protobuf:"bytes,2,rep,name=params"`
	Disabled bool                  `json:"disabled,omitempty" protobuf:"varint,3,opt,name=disabled"`
}

// TODO: describe this type
// TODO: describe members of this type
type ResourceActionParam struct {
	Name    string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Value   string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
	Type    string `json:"type,omitempty" protobuf:"bytes,3,opt,name=type"`
	Default string `json:"default,omitempty" protobuf:"bytes,4,opt,name=default"`
}
