package validate

import (
	"strings"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllValidators(t *testing.T) {
	for _, validator := range validators.TestRegistry.All() {
		validator := validator
		t.Run(validator.Name(), func(t *testing.T) {
			t.Parallel()
			// testing the succeeding candidates
			succeedingMetaBundles := validator.SucceedingCandidates()
			for _, mb := range succeedingMetaBundles {
				res := validator.Run(mb)
				require.False(t, res.IsError())
				assert.True(t, res.IsSuccess())
			}

			// (optional) testing the failing candidates
			failingMetaBundles := validator.FailingCandidates()
			for _, mb := range failingMetaBundles {
				res := validator.Run(mb)
				require.False(t, res.IsError())
				assert.False(t, res.IsSuccess())
			}
		})
	}
}

func TestFilterDisabledValidators(t *testing.T) {
	n_validators := validators.Registry.Len()

	cases := []struct {
		name     string
		disabled []string
	}{
		{
			name:     "all_enabled",
			disabled: []string{},
		},
		{
			name:     "disable_default_channel",
			disabled: []string{"AM0001"},
		},
		{
			name:     "disable_two",
			disabled: []string{"AM0001", "AM0002"},
		},
	}
	for _, tc := range cases {
		tc := tc // pin
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			filter, err := NewFilter(strings.Join(tc.disabled, ","), "")
			require.NoError(t, err)

			n_enabled := len(filter.GetValidators())
			n_disabled := len(tc.disabled)
			require.Equal(t, n_enabled+n_disabled, n_validators)
		})
	}
}

func TestFilterEnabledValidators(t *testing.T) {
	cases := []struct {
		name    string
		enabled []string
	}{
		{
			name:    "enable_default_channel",
			enabled: []string{"AM0001"},
		},
		{
			name:    "enable_two",
			enabled: []string{"AM0001", "AM0002"},
		},
	}
	for _, tc := range cases {
		tc := tc // pin
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			filter, err := NewFilter("", strings.Join(tc.enabled, ","))
			require.NoError(t, err)
			require.Equal(t, len(filter.GetValidators()), len(tc.enabled))
		})
	}
}

func TestEmptyFilterAllEnabled(t *testing.T) {
	t.Parallel()
	filter, err := NewFilter("", "")
	require.NoError(t, err)
	require.Equal(t, len(filter.GetValidators()), validators.Registry.Len())
}

func TestFilterError(t *testing.T) {
	cases := []struct {
		name     string
		enabled  []string
		disabled []string
	}{
		{
			name:     "mutually_exclusive",
			enabled:  []string{"AM0001"},
			disabled: []string{"AM0001"},
		},
		{
			name:     "enabled_dont_exist",
			enabled:  []string{"invalid"},
			disabled: []string{},
		},
		{
			name:     "disabled_dont_exist",
			enabled:  []string{},
			disabled: []string{"invalid"},
		},
	}
	for _, tc := range cases {
		tc := tc // pin
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			disabled := strings.Join(tc.disabled, ",")
			enabled := strings.Join(tc.enabled, ",")
			filter, err := NewFilter(disabled, enabled)
			require.Error(t, err)
			require.Nil(t, filter)
		})
	}
}
