// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"context"
	"sync"

	"github.com/51st-state/api/pkg/rbac"
)

type FakeControl struct {
	GetRoleRulesStub        func(ctx context.Context, roleID rbac.RoleID) (rbac.RoleRules, error)
	getRoleRulesMutex       sync.RWMutex
	getRoleRulesArgsForCall []struct {
		ctx    context.Context
		roleID rbac.RoleID
	}
	getRoleRulesReturns struct {
		result1 rbac.RoleRules
		result2 error
	}
	getRoleRulesReturnsOnCall map[int]struct {
		result1 rbac.RoleRules
		result2 error
	}
	SetRoleRulesStub        func(ctx context.Context, roleID rbac.RoleID, rules rbac.RoleRules) error
	setRoleRulesMutex       sync.RWMutex
	setRoleRulesArgsForCall []struct {
		ctx    context.Context
		roleID rbac.RoleID
		rules  rbac.RoleRules
	}
	setRoleRulesReturns struct {
		result1 error
	}
	setRoleRulesReturnsOnCall map[int]struct {
		result1 error
	}
	GetSubjectRolesStub        func(ctx context.Context, subjectID rbac.SubjectID) (rbac.SubjectRoles, error)
	getSubjectRolesMutex       sync.RWMutex
	getSubjectRolesArgsForCall []struct {
		ctx       context.Context
		subjectID rbac.SubjectID
	}
	getSubjectRolesReturns struct {
		result1 rbac.SubjectRoles
		result2 error
	}
	getSubjectRolesReturnsOnCall map[int]struct {
		result1 rbac.SubjectRoles
		result2 error
	}
	SetSubjectRolesStub        func(ctx context.Context, subjectID rbac.SubjectID, roles rbac.SubjectRoles) error
	setSubjectRolesMutex       sync.RWMutex
	setSubjectRolesArgsForCall []struct {
		ctx       context.Context
		subjectID rbac.SubjectID
		roles     rbac.SubjectRoles
	}
	setSubjectRolesReturns struct {
		result1 error
	}
	setSubjectRolesReturnsOnCall map[int]struct {
		result1 error
	}
	IsSubjectAllowedStub        func(ctx context.Context, subjectID rbac.SubjectID, rule rbac.Rule) error
	isSubjectAllowedMutex       sync.RWMutex
	isSubjectAllowedArgsForCall []struct {
		ctx       context.Context
		subjectID rbac.SubjectID
		rule      rbac.Rule
	}
	isSubjectAllowedReturns struct {
		result1 error
	}
	isSubjectAllowedReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeControl) GetRoleRules(ctx context.Context, roleID rbac.RoleID) (rbac.RoleRules, error) {
	fake.getRoleRulesMutex.Lock()
	ret, specificReturn := fake.getRoleRulesReturnsOnCall[len(fake.getRoleRulesArgsForCall)]
	fake.getRoleRulesArgsForCall = append(fake.getRoleRulesArgsForCall, struct {
		ctx    context.Context
		roleID rbac.RoleID
	}{ctx, roleID})
	fake.recordInvocation("GetRoleRules", []interface{}{ctx, roleID})
	fake.getRoleRulesMutex.Unlock()
	if fake.GetRoleRulesStub != nil {
		return fake.GetRoleRulesStub(ctx, roleID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getRoleRulesReturns.result1, fake.getRoleRulesReturns.result2
}

func (fake *FakeControl) GetRoleRulesCallCount() int {
	fake.getRoleRulesMutex.RLock()
	defer fake.getRoleRulesMutex.RUnlock()
	return len(fake.getRoleRulesArgsForCall)
}

func (fake *FakeControl) GetRoleRulesArgsForCall(i int) (context.Context, rbac.RoleID) {
	fake.getRoleRulesMutex.RLock()
	defer fake.getRoleRulesMutex.RUnlock()
	return fake.getRoleRulesArgsForCall[i].ctx, fake.getRoleRulesArgsForCall[i].roleID
}

func (fake *FakeControl) GetRoleRulesReturns(result1 rbac.RoleRules, result2 error) {
	fake.GetRoleRulesStub = nil
	fake.getRoleRulesReturns = struct {
		result1 rbac.RoleRules
		result2 error
	}{result1, result2}
}

func (fake *FakeControl) GetRoleRulesReturnsOnCall(i int, result1 rbac.RoleRules, result2 error) {
	fake.GetRoleRulesStub = nil
	if fake.getRoleRulesReturnsOnCall == nil {
		fake.getRoleRulesReturnsOnCall = make(map[int]struct {
			result1 rbac.RoleRules
			result2 error
		})
	}
	fake.getRoleRulesReturnsOnCall[i] = struct {
		result1 rbac.RoleRules
		result2 error
	}{result1, result2}
}

func (fake *FakeControl) SetRoleRules(ctx context.Context, roleID rbac.RoleID, rules rbac.RoleRules) error {
	fake.setRoleRulesMutex.Lock()
	ret, specificReturn := fake.setRoleRulesReturnsOnCall[len(fake.setRoleRulesArgsForCall)]
	fake.setRoleRulesArgsForCall = append(fake.setRoleRulesArgsForCall, struct {
		ctx    context.Context
		roleID rbac.RoleID
		rules  rbac.RoleRules
	}{ctx, roleID, rules})
	fake.recordInvocation("SetRoleRules", []interface{}{ctx, roleID, rules})
	fake.setRoleRulesMutex.Unlock()
	if fake.SetRoleRulesStub != nil {
		return fake.SetRoleRulesStub(ctx, roleID, rules)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setRoleRulesReturns.result1
}

func (fake *FakeControl) SetRoleRulesCallCount() int {
	fake.setRoleRulesMutex.RLock()
	defer fake.setRoleRulesMutex.RUnlock()
	return len(fake.setRoleRulesArgsForCall)
}

func (fake *FakeControl) SetRoleRulesArgsForCall(i int) (context.Context, rbac.RoleID, rbac.RoleRules) {
	fake.setRoleRulesMutex.RLock()
	defer fake.setRoleRulesMutex.RUnlock()
	return fake.setRoleRulesArgsForCall[i].ctx, fake.setRoleRulesArgsForCall[i].roleID, fake.setRoleRulesArgsForCall[i].rules
}

func (fake *FakeControl) SetRoleRulesReturns(result1 error) {
	fake.SetRoleRulesStub = nil
	fake.setRoleRulesReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeControl) SetRoleRulesReturnsOnCall(i int, result1 error) {
	fake.SetRoleRulesStub = nil
	if fake.setRoleRulesReturnsOnCall == nil {
		fake.setRoleRulesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setRoleRulesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeControl) GetSubjectRoles(ctx context.Context, subjectID rbac.SubjectID) (rbac.SubjectRoles, error) {
	fake.getSubjectRolesMutex.Lock()
	ret, specificReturn := fake.getSubjectRolesReturnsOnCall[len(fake.getSubjectRolesArgsForCall)]
	fake.getSubjectRolesArgsForCall = append(fake.getSubjectRolesArgsForCall, struct {
		ctx       context.Context
		subjectID rbac.SubjectID
	}{ctx, subjectID})
	fake.recordInvocation("GetSubjectRoles", []interface{}{ctx, subjectID})
	fake.getSubjectRolesMutex.Unlock()
	if fake.GetSubjectRolesStub != nil {
		return fake.GetSubjectRolesStub(ctx, subjectID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getSubjectRolesReturns.result1, fake.getSubjectRolesReturns.result2
}

func (fake *FakeControl) GetSubjectRolesCallCount() int {
	fake.getSubjectRolesMutex.RLock()
	defer fake.getSubjectRolesMutex.RUnlock()
	return len(fake.getSubjectRolesArgsForCall)
}

func (fake *FakeControl) GetSubjectRolesArgsForCall(i int) (context.Context, rbac.SubjectID) {
	fake.getSubjectRolesMutex.RLock()
	defer fake.getSubjectRolesMutex.RUnlock()
	return fake.getSubjectRolesArgsForCall[i].ctx, fake.getSubjectRolesArgsForCall[i].subjectID
}

func (fake *FakeControl) GetSubjectRolesReturns(result1 rbac.SubjectRoles, result2 error) {
	fake.GetSubjectRolesStub = nil
	fake.getSubjectRolesReturns = struct {
		result1 rbac.SubjectRoles
		result2 error
	}{result1, result2}
}

func (fake *FakeControl) GetSubjectRolesReturnsOnCall(i int, result1 rbac.SubjectRoles, result2 error) {
	fake.GetSubjectRolesStub = nil
	if fake.getSubjectRolesReturnsOnCall == nil {
		fake.getSubjectRolesReturnsOnCall = make(map[int]struct {
			result1 rbac.SubjectRoles
			result2 error
		})
	}
	fake.getSubjectRolesReturnsOnCall[i] = struct {
		result1 rbac.SubjectRoles
		result2 error
	}{result1, result2}
}

func (fake *FakeControl) SetSubjectRoles(ctx context.Context, subjectID rbac.SubjectID, roles rbac.SubjectRoles) error {
	fake.setSubjectRolesMutex.Lock()
	ret, specificReturn := fake.setSubjectRolesReturnsOnCall[len(fake.setSubjectRolesArgsForCall)]
	fake.setSubjectRolesArgsForCall = append(fake.setSubjectRolesArgsForCall, struct {
		ctx       context.Context
		subjectID rbac.SubjectID
		roles     rbac.SubjectRoles
	}{ctx, subjectID, roles})
	fake.recordInvocation("SetSubjectRoles", []interface{}{ctx, subjectID, roles})
	fake.setSubjectRolesMutex.Unlock()
	if fake.SetSubjectRolesStub != nil {
		return fake.SetSubjectRolesStub(ctx, subjectID, roles)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setSubjectRolesReturns.result1
}

func (fake *FakeControl) SetSubjectRolesCallCount() int {
	fake.setSubjectRolesMutex.RLock()
	defer fake.setSubjectRolesMutex.RUnlock()
	return len(fake.setSubjectRolesArgsForCall)
}

func (fake *FakeControl) SetSubjectRolesArgsForCall(i int) (context.Context, rbac.SubjectID, rbac.SubjectRoles) {
	fake.setSubjectRolesMutex.RLock()
	defer fake.setSubjectRolesMutex.RUnlock()
	return fake.setSubjectRolesArgsForCall[i].ctx, fake.setSubjectRolesArgsForCall[i].subjectID, fake.setSubjectRolesArgsForCall[i].roles
}

func (fake *FakeControl) SetSubjectRolesReturns(result1 error) {
	fake.SetSubjectRolesStub = nil
	fake.setSubjectRolesReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeControl) SetSubjectRolesReturnsOnCall(i int, result1 error) {
	fake.SetSubjectRolesStub = nil
	if fake.setSubjectRolesReturnsOnCall == nil {
		fake.setSubjectRolesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setSubjectRolesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeControl) IsSubjectAllowed(ctx context.Context, subjectID rbac.SubjectID, rule rbac.Rule) error {
	fake.isSubjectAllowedMutex.Lock()
	ret, specificReturn := fake.isSubjectAllowedReturnsOnCall[len(fake.isSubjectAllowedArgsForCall)]
	fake.isSubjectAllowedArgsForCall = append(fake.isSubjectAllowedArgsForCall, struct {
		ctx       context.Context
		subjectID rbac.SubjectID
		rule      rbac.Rule
	}{ctx, subjectID, rule})
	fake.recordInvocation("IsSubjectAllowed", []interface{}{ctx, subjectID, rule})
	fake.isSubjectAllowedMutex.Unlock()
	if fake.IsSubjectAllowedStub != nil {
		return fake.IsSubjectAllowedStub(ctx, subjectID, rule)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.isSubjectAllowedReturns.result1
}

func (fake *FakeControl) IsSubjectAllowedCallCount() int {
	fake.isSubjectAllowedMutex.RLock()
	defer fake.isSubjectAllowedMutex.RUnlock()
	return len(fake.isSubjectAllowedArgsForCall)
}

func (fake *FakeControl) IsSubjectAllowedArgsForCall(i int) (context.Context, rbac.SubjectID, rbac.Rule) {
	fake.isSubjectAllowedMutex.RLock()
	defer fake.isSubjectAllowedMutex.RUnlock()
	return fake.isSubjectAllowedArgsForCall[i].ctx, fake.isSubjectAllowedArgsForCall[i].subjectID, fake.isSubjectAllowedArgsForCall[i].rule
}

func (fake *FakeControl) IsSubjectAllowedReturns(result1 error) {
	fake.IsSubjectAllowedStub = nil
	fake.isSubjectAllowedReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeControl) IsSubjectAllowedReturnsOnCall(i int, result1 error) {
	fake.IsSubjectAllowedStub = nil
	if fake.isSubjectAllowedReturnsOnCall == nil {
		fake.isSubjectAllowedReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.isSubjectAllowedReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeControl) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getRoleRulesMutex.RLock()
	defer fake.getRoleRulesMutex.RUnlock()
	fake.setRoleRulesMutex.RLock()
	defer fake.setRoleRulesMutex.RUnlock()
	fake.getSubjectRolesMutex.RLock()
	defer fake.getSubjectRolesMutex.RUnlock()
	fake.setSubjectRolesMutex.RLock()
	defer fake.setSubjectRolesMutex.RUnlock()
	fake.isSubjectAllowedMutex.RLock()
	defer fake.isSubjectAllowedMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeControl) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ rbac.Control = new(FakeControl)