/*
Copyright 2021 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package services

import (
	"github.com/google/uuid"
	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	"github.com/gravitational/teleport"
	"github.com/gravitational/teleport/api/constants"
	apidefaults "github.com/gravitational/teleport/api/defaults"
	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/lib/modules"
)

// NewPresetEditorRole returns a new pre-defined role for cluster
// editors who can edit cluster configuration resources.
func NewPresetEditorRole() types.Role {
	enterprise := modules.GetModules().BuildType() == modules.BuildEnterprise
	role := &types.RoleV6{
		Kind:    types.KindRole,
		Version: types.V7,
		Metadata: types.Metadata{
			Name:        teleport.PresetEditorRoleName,
			Namespace:   apidefaults.Namespace,
			Description: "Edit cluster configuration",
		},
		Spec: types.RoleSpecV6{
			Options: types.RoleOptions{
				CertificateFormat: constants.CertificateFormatStandard,
				MaxSessionTTL:     types.NewDuration(apidefaults.MaxCertDuration),
				PortForwarding:    types.NewBoolOption(true),
				ForwardAgent:      types.NewBool(true),
				BPF:               apidefaults.EnhancedEvents(),
				RecordSession: &types.RecordSession{
					Desktop: types.NewBoolOption(false),
				},
			},
			Allow: types.RoleConditions{
				Namespaces: []string{apidefaults.Namespace},
				Rules: []types.Rule{
					types.NewRule(types.KindUser, RW()),
					types.NewRule(types.KindRole, RW()),
					types.NewRule(types.KindOIDC, RW()),
					types.NewRule(types.KindSAML, RW()),
					types.NewRule(types.KindGithub, RW()),
					types.NewRule(types.KindOIDCRequest, RW()),
					types.NewRule(types.KindSAMLRequest, RW()),
					types.NewRule(types.KindGithubRequest, RW()),
					types.NewRule(types.KindClusterAuditConfig, RW()),
					types.NewRule(types.KindClusterAuthPreference, RW()),
					types.NewRule(types.KindAuthConnector, RW()),
					types.NewRule(types.KindClusterName, RW()),
					types.NewRule(types.KindClusterNetworkingConfig, RW()),
					types.NewRule(types.KindSessionRecordingConfig, RW()),
					types.NewRule(types.KindUIConfig, RW()),
					types.NewRule(types.KindTrustedCluster, RW()),
					types.NewRule(types.KindRemoteCluster, RW()),
					types.NewRule(types.KindToken, RW()),
					types.NewRule(types.KindConnectionDiagnostic, RW()),
					types.NewRule(types.KindDatabase, RW()),
					types.NewRule(types.KindDatabaseCertificate, RW()),
					types.NewRule(types.KindInstaller, RW()),
					types.NewRule(types.KindDevice, append(RW(), types.VerbCreateEnrollToken, types.VerbEnroll)),
					types.NewRule(types.KindDatabaseService, RO()),
					types.NewRule(types.KindInstance, RO()),
					types.NewRule(types.KindLoginRule, RW()),
					types.NewRule(types.KindSAMLIdPServiceProvider, RW()),
					types.NewRule(types.KindUserGroup, RW()),
					types.NewRule(types.KindPlugin, RW()),
					types.NewRule(types.KindOktaImportRule, RW()),
					types.NewRule(types.KindOktaAssignment, RW()),
					types.NewRule(types.KindAssistant, append(RW(), types.VerbUse)),
					types.NewRule(types.KindLock, RW()),
					types.NewRule(types.KindIntegration, append(RW(), types.VerbUse)),
					types.NewRule(types.KindBilling, RW()),
					types.NewRule(types.KindClusterAlert, RW()),
					// Please see defaultAllowRules when adding a new rule.
				},
				// By default, allow editors to approve any user group access requests.
				ReviewRequests: defaultAllowAccessReviewConditions(enterprise)[teleport.PresetEditorRoleName],
			},
		},
	}
	return role
}

// NewPresetAccessRole creates a role for users who are allowed to initiate
// interactive sessions.
func NewPresetAccessRole() types.Role {
	enterprise := modules.GetModules().BuildType() == modules.BuildEnterprise
	role := &types.RoleV6{
		Kind:    types.KindRole,
		Version: types.V7,
		Metadata: types.Metadata{
			Name:        teleport.PresetAccessRoleName,
			Namespace:   apidefaults.Namespace,
			Description: "Access cluster resources",
		},
		Spec: types.RoleSpecV6{
			Options: types.RoleOptions{
				CertificateFormat: constants.CertificateFormatStandard,
				MaxSessionTTL:     types.NewDuration(apidefaults.MaxCertDuration),
				PortForwarding:    types.NewBoolOption(true),
				ForwardAgent:      types.NewBool(true),
				BPF:               apidefaults.EnhancedEvents(),
				RecordSession:     &types.RecordSession{Desktop: types.NewBoolOption(true)},
			},
			Allow: types.RoleConditions{
				Namespaces:            []string{apidefaults.Namespace},
				NodeLabels:            types.Labels{types.Wildcard: []string{types.Wildcard}},
				AppLabels:             types.Labels{types.Wildcard: []string{types.Wildcard}},
				KubernetesLabels:      types.Labels{types.Wildcard: []string{types.Wildcard}},
				WindowsDesktopLabels:  types.Labels{types.Wildcard: []string{types.Wildcard}},
				DatabaseLabels:        types.Labels{types.Wildcard: []string{types.Wildcard}},
				DatabaseServiceLabels: types.Labels{types.Wildcard: []string{types.Wildcard}},
				DatabaseNames:         []string{teleport.TraitInternalDBNamesVariable},
				DatabaseUsers:         []string{teleport.TraitInternalDBUsersVariable},
				DatabaseRoles:         []string{teleport.TraitInternalDBRolesVariable},
				KubernetesResources: []types.KubernetesResource{
					{
						Kind:      types.Wildcard,
						Namespace: types.Wildcard,
						Name:      types.Wildcard,
					},
				},
				Rules: []types.Rule{
					types.NewRule(types.KindEvent, RO()),
					{
						Resources: []string{types.KindSession},
						Verbs:     []string{types.VerbRead, types.VerbList},
						Where:     "contains(session.participants, user.metadata.name)",
					},
					types.NewRule(types.KindInstance, RO()),
					types.NewRule(types.KindAssistant, append(RW(), types.VerbUse)),
					// Please see defaultAllowRules when adding a new rule.
				},
				// By default, allow users with the access role to request any user group.
				Request: defaultAllowAccessRequestConditions(enterprise)[teleport.PresetAccessRoleName],
			},
		},
	}
	role.SetLogins(types.Allow, []string{teleport.TraitInternalLoginsVariable})
	role.SetWindowsLogins(types.Allow, []string{teleport.TraitInternalWindowsLoginsVariable})
	role.SetKubeUsers(types.Allow, []string{teleport.TraitInternalKubeUsersVariable})
	role.SetKubeGroups(types.Allow, []string{teleport.TraitInternalKubeGroupsVariable})
	role.SetAWSRoleARNs(types.Allow, []string{teleport.TraitInternalAWSRoleARNs})
	role.SetAzureIdentities(types.Allow, []string{teleport.TraitInternalAzureIdentities})
	role.SetGCPServiceAccounts(types.Allow, []string{teleport.TraitInternalGCPServiceAccounts})
	return role
}

// NewPresetAuditorRole returns a new pre-defined role for cluster
// auditor - someone who can review cluster events and replay sessions,
// but can't initiate interactive sessions or modify configuration.
func NewPresetAuditorRole() types.Role {
	role := &types.RoleV6{
		Kind:    types.KindRole,
		Version: types.V7,
		Metadata: types.Metadata{
			Name:        teleport.PresetAuditorRoleName,
			Namespace:   apidefaults.Namespace,
			Description: "Review cluster events and replay sessions",
		},
		Spec: types.RoleSpecV6{
			Options: types.RoleOptions{
				CertificateFormat: constants.CertificateFormatStandard,
				MaxSessionTTL:     types.NewDuration(apidefaults.MaxCertDuration),
				RecordSession: &types.RecordSession{
					Desktop: types.NewBoolOption(false),
				},
			},
			Allow: types.RoleConditions{
				Namespaces: []string{apidefaults.Namespace},
				Rules: []types.Rule{
					types.NewRule(types.KindSession, RO()),
					types.NewRule(types.KindEvent, RO()),
					types.NewRule(types.KindSessionTracker, RO()),
					types.NewRule(types.KindClusterAlert, RO()),
					// Please see defaultAllowRules when adding a new rule.
				},
			},
		},
	}
	role.SetLogins(types.Allow, []string{"no-login-" + uuid.New().String()})
	return role
}

// NewPresetGroupAccessRole returns a new pre-defined role for group access -
// a role used for requesting and reviewing user group access.
func NewPresetGroupAccessRole() types.Role {
	role := &types.RoleV6{
		Kind:    types.KindRole,
		Version: types.V6,
		Metadata: types.Metadata{
			Name:        teleport.PresetGroupAccessRoleName,
			Namespace:   apidefaults.Namespace,
			Description: "Have access to all user groups",
		},
		Spec: types.RoleSpecV6{
			Allow: types.RoleConditions{
				Namespaces: []string{apidefaults.Namespace},
				GroupLabels: types.Labels{
					types.Wildcard: []string{types.Wildcard},
				},
				Rules: []types.Rule{
					types.NewRule(types.KindUserGroup, RO()),
					// Please see defaultAllowRules when adding a new rule.
				},
			},
		},
	}
	role.SetLogins(types.Allow, []string{"no-login-" + uuid.New().String()})
	return role
}

// defaultAllowRules has the Allow rules that should be set as default when
// they were not explicitly defined. This is used to update the current cluster
// roles when deploying a new resource. It will also update all existing roles
// on auth server restart. Rules defined in preset template should be
// exactly the same rule when added here.
func defaultAllowRules() map[string][]types.Rule {
	return map[string][]types.Rule{
		teleport.PresetAuditorRoleName: {
			types.NewRule(types.KindSessionTracker, RO()),
		},
		teleport.PresetEditorRoleName: {
			types.NewRule(types.KindConnectionDiagnostic, RW()),
			types.NewRule(types.KindDatabase, RW()),
			types.NewRule(types.KindDatabaseService, RO()),
			types.NewRule(types.KindLoginRule, RW()),
			types.NewRule(types.KindPlugin, RW()),
			types.NewRule(types.KindSAMLIdPServiceProvider, RW()),
			types.NewRule(types.KindOktaImportRule, RW()),
			types.NewRule(types.KindOktaAssignment, RW()),
			types.NewRule(types.KindDevice, append(RW(), types.VerbCreateEnrollToken, types.VerbEnroll)),
			types.NewRule(types.KindLock, RW()),
			types.NewRule(types.KindIntegration, append(RW(), types.VerbUse)),
			types.NewRule(types.KindBilling, RW()),
			types.NewRule(types.KindAssistant, append(RW(), types.VerbUse)),
		},
		teleport.PresetAccessRoleName: {
			// Allow assist access to access role. This role only allow access
			// to the assist console, not any other cluster resources.
			types.NewRule(types.KindAssistant, append(RW(), types.VerbUse)),
		},
	}
}

// defaultAllowLabels has the Allow labels that should be set as default when they were not explicitly defined.
// This is used to update existing builtin preset roles with new permissions during cluster upgrades.
// The following Labels are supported:
// - DatabaseServiceLabels (db_service_labels)
func defaultAllowLabels() map[string]types.RoleConditions {
	return map[string]types.RoleConditions{
		teleport.PresetAccessRoleName: {
			DatabaseServiceLabels: types.Labels{types.Wildcard: []string{types.Wildcard}},
			DatabaseRoles:         []string{teleport.TraitInternalDBRolesVariable},
		},
	}
}

// defaultAllowAccessRequestConditions has the access request conditions that should be set as default when they were
// not explicitly defined.
func defaultAllowAccessRequestConditions(enterprise bool) map[string]*types.AccessRequestConditions {
	if enterprise {
		return map[string]*types.AccessRequestConditions{
			teleport.PresetAccessRoleName: {
				SearchAsRoles: []string{
					teleport.PresetGroupAccessRoleName,
				},
			},
		}
	}

	return map[string]*types.AccessRequestConditions{}
}

// defaultAllowAccessReviewConditions has the access review conditions that should be set as default when they were
// not explicitly defined.
func defaultAllowAccessReviewConditions(enterprise bool) map[string]*types.AccessReviewConditions {
	if enterprise {
		return map[string]*types.AccessReviewConditions{
			teleport.PresetEditorRoleName: {
				PreviewAsRoles: []string{
					teleport.PresetGroupAccessRoleName,
				},
				Roles: []string{
					teleport.PresetGroupAccessRoleName,
				},
			},
		}
	}

	return map[string]*types.AccessReviewConditions{}
}

// AddRoleDefaults adds default role attributes to a preset role.
// Only attributes whose resources are not already defined (either allowing or denying) are added.
func AddRoleDefaults(role types.Role) (types.Role, error) {
	changed := false

	// Resource Rules
	defaultRules, ok := defaultAllowRules()[role.GetName()]
	if ok {
		existingRules := append(role.GetRules(types.Allow), role.GetRules(types.Deny)...)

		for _, defaultRule := range defaultRules {
			if resourceBelongsToRules(existingRules, defaultRule.Resources) {
				continue
			}

			log.Debugf("Adding default allow rule %v for role %q", defaultRule, role.GetName())
			rules := role.GetRules(types.Allow)
			rules = append(rules, defaultRule)
			role.SetRules(types.Allow, rules)
			changed = true
		}
	}

	// Labels
	defaultLabels, ok := defaultAllowLabels()[role.GetName()]
	if ok {
		if unset, err := labelMatchersUnset(role, types.KindDatabaseService); err != nil {
			return nil, trace.Wrap(err)
		} else if unset && len(defaultLabels.DatabaseServiceLabels) > 0 {
			role.SetLabelMatchers(types.Allow, types.KindDatabaseService, types.LabelMatchers{
				Labels: defaultLabels.DatabaseServiceLabels,
			})
			changed = true
		}
		if len(defaultLabels.DatabaseRoles) > 0 && len(role.GetDatabaseRoles(types.Allow)) == 0 {
			role.SetDatabaseRoles(types.Allow, defaultLabels.DatabaseRoles)
			changed = true
		}
	}

	enterprise := modules.GetModules().BuildType() == modules.BuildEnterprise

	if role.GetAccessRequestConditions(types.Allow).IsEmpty() {
		arc := defaultAllowAccessRequestConditions(enterprise)[role.GetName()]
		if arc != nil {
			role.SetAccessRequestConditions(types.Allow, *arc)
			changed = true
		}
	}

	if role.GetAccessReviewConditions(types.Allow).IsEmpty() {
		arc := defaultAllowAccessReviewConditions(enterprise)[role.GetName()]
		if arc != nil {
			role.SetAccessReviewConditions(types.Allow, *arc)
			changed = true
		}
	}

	if !changed {
		return nil, trace.AlreadyExists("no change")
	}

	return role, nil
}

func labelMatchersUnset(role types.Role, kind string) (bool, error) {
	for _, cond := range []types.RoleConditionType{types.Allow, types.Deny} {
		labelMatchers, err := role.GetLabelMatchers(cond, kind)
		if err != nil {
			return false, trace.Wrap(err)
		}
		if !labelMatchers.Empty() {
			return false, nil
		}
	}
	return true, nil
}

func resourceBelongsToRules(rules []types.Rule, resources []string) bool {
	for _, rule := range rules {
		for _, ruleResource := range rule.Resources {
			if slices.Contains(resources, ruleResource) {
				return true
			}
		}
	}

	return false
}
