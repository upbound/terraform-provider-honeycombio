package modifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestDatasetDeprecation(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		allowAnyDataset bool
		config          types.String
		state           types.String
		plan            types.String
		expectedPlan    types.String
		expectReplace   bool
	}{
		"unset dataset, another attribute changed (plan marked unknown)": {
			allowAnyDataset: true,
			config:          types.StringNull(),
			state:           types.StringNull(),
			plan:            types.StringUnknown(),
			expectedPlan:    types.StringNull(),
			expectReplace:   false,
		},
		"unset dataset, nothing changed": {
			allowAnyDataset: true,
			config:          types.StringNull(),
			state:           types.StringNull(),
			plan:            types.StringNull(),
			expectedPlan:    types.StringNull(),
			expectReplace:   false,
		},
		"dataset removed from config, allowAnyDataset keeps state": {
			allowAnyDataset: true,
			config:          types.StringNull(),
			state:           types.StringValue("foo"),
			plan:            types.StringUnknown(),
			expectedPlan:    types.StringValue("foo"),
			expectReplace:   false,
		},
		"dataset removed from config, strict mode replaces": {
			allowAnyDataset: false,
			config:          types.StringNull(),
			state:           types.StringValue("foo"),
			plan:            types.StringUnknown(),
			expectedPlan:    types.StringNull(),
			expectReplace:   true,
		},
		"dataset changed": {
			allowAnyDataset: true,
			config:          types.StringValue("bar"),
			state:           types.StringValue("foo"),
			plan:            types.StringValue("bar"),
			expectedPlan:    types.StringValue("bar"),
			expectReplace:   true,
		},
		"dataset unchanged": {
			allowAnyDataset: true,
			config:          types.StringValue("foo"),
			state:           types.StringValue("foo"),
			plan:            types.StringValue("foo"),
			expectedPlan:    types.StringValue("foo"),
			expectReplace:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.StringRequest{
				ConfigValue: tc.config,
				StateValue:  tc.state,
				PlanValue:   tc.plan,
				Plan:        tfsdk.Plan{Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{})},
			}
			resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}

			DatasetDeprecation(tc.allowAnyDataset).PlanModifyString(context.Background(), req, resp)

			if !resp.PlanValue.Equal(tc.expectedPlan) {
				t.Errorf("PlanValue: got %s, want %s", resp.PlanValue, tc.expectedPlan)
			}
			if resp.RequiresReplace != tc.expectReplace {
				t.Errorf("RequiresReplace: got %v, want %v", resp.RequiresReplace, tc.expectReplace)
			}
		})
	}
}
