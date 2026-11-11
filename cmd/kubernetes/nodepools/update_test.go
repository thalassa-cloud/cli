package nodepools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

func TestBuildNodePoolUpdateRequest(t *testing.T) {
	strategyManual := kubernetes.KubernetesNodePoolUpgradeStrategyManual

	tests := []struct {
		name      string
		current   *kubernetes.KubernetesNodePool
		overrides NodePoolUpdateOverrides
		wantErr   string
		assertReq func(t *testing.T, current *kubernetes.KubernetesNodePool, req kubernetes.UpdateKubernetesNodePool)
	}{
		{
			name:    "copies the current node pool",
			current: sampleNodePool(),
			assertReq: func(t *testing.T, current *kubernetes.KubernetesNodePool, req kubernetes.UpdateKubernetesNodePool) {
				assert.Equal(t, current.Description, req.Description)
				assert.Equal(t, current.Labels, req.Labels)
				assert.Equal(t, current.Annotations, req.Annotations)
				assert.Equal(t, current.MachineType.Name, req.MachineType)
				require.NotNil(t, req.Replicas)
				assert.Equal(t, current.Replicas, *req.Replicas)
				require.NotNil(t, req.MinReplicas)
				assert.Equal(t, current.MinReplicas, *req.MinReplicas)
				require.NotNil(t, req.MaxReplicas)
				assert.Equal(t, current.MaxReplicas, *req.MaxReplicas)
				require.NotNil(t, req.KubernetesVersionIdentity)
				assert.Equal(t, current.KubernetesVersion.Identity, *req.KubernetesVersionIdentity)
				assert.Equal(t, current.AvailabilityZone, req.AvailabilityZone)
				require.NotNil(t, req.UpgradeStrategy)
				assert.Equal(t, current.UpgradeStrategy, *req.UpgradeStrategy)
				require.NotNil(t, req.EnableAutoscaling)
				assert.Equal(t, current.EnableAutoscaling, *req.EnableAutoscaling)
				require.NotNil(t, req.EnableAutoHealing)
				assert.Equal(t, current.EnableAutoHealing, *req.EnableAutoHealing)
				assert.Equal(t, current.ManageNodeAllocatable, req.ManageNodeAllocatable)
				require.NotNil(t, req.NodeSettings)
				assert.Equal(t, current.NodeSettings.Labels, req.NodeSettings.Labels)
				assert.Equal(t, current.NodeSettings.Annotations, req.NodeSettings.Annotations)
				assert.Equal(t, current.NodeSettings.Taints, req.NodeSettings.Taints)
				assert.Equal(t, []string{"sg-1", "sg-2"}, req.SecurityGroupAttachments)
			},
		},
		{
			name:    "overrides the version and keeps the rest",
			current: sampleNodePool(),
			overrides: NodePoolUpdateOverrides{
				KubernetesVersionIdentity: stringPtr("ver-2"),
			},
			assertReq: func(t *testing.T, current *kubernetes.KubernetesNodePool, req kubernetes.UpdateKubernetesNodePool) {
				require.NotNil(t, req.KubernetesVersionIdentity)
				assert.Equal(t, "ver-2", *req.KubernetesVersionIdentity)
				assert.Equal(t, "pgp-medium", req.MachineType)
				require.NotNil(t, req.EnableAutoscaling)
				assert.True(t, *req.EnableAutoscaling)
				require.NotNil(t, req.Replicas)
				assert.Equal(t, 2, *req.Replicas)
				assert.Equal(t, []string{"sg-1", "sg-2"}, req.SecurityGroupAttachments)
				require.NotNil(t, req.NodeSettings)
				assert.Equal(t, current.NodeSettings.Taints, req.NodeSettings.Taints)
				assert.Equal(t, current.NodeSettings.Labels, req.NodeSettings.Labels)
			},
		},
		{
			name:    "overrides size and strategy",
			current: sampleNodePool(),
			overrides: NodePoolUpdateOverrides{
				MachineType:       stringPtr("pgp-large"),
				Replicas:          intPtr(4),
				MinReplicas:       intPtr(2),
				MaxReplicas:       intPtr(8),
				EnableAutoscaling: boolPtr(false),
				EnableAutoHealing: boolPtr(false),
				UpgradeStrategy:   &strategyManual,
			},
			assertReq: func(t *testing.T, _ *kubernetes.KubernetesNodePool, req kubernetes.UpdateKubernetesNodePool) {
				assert.Equal(t, "pgp-large", req.MachineType)
				require.NotNil(t, req.Replicas)
				assert.Equal(t, 4, *req.Replicas)
				require.NotNil(t, req.MinReplicas)
				assert.Equal(t, 2, *req.MinReplicas)
				require.NotNil(t, req.MaxReplicas)
				assert.Equal(t, 8, *req.MaxReplicas)
				require.NotNil(t, req.EnableAutoscaling)
				assert.False(t, *req.EnableAutoscaling)
				require.NotNil(t, req.EnableAutoHealing)
				assert.False(t, *req.EnableAutoHealing)
				require.NotNil(t, req.UpgradeStrategy)
				assert.Equal(t, strategyManual, *req.UpgradeStrategy)
				require.NotNil(t, req.KubernetesVersionIdentity)
				assert.Equal(t, "ver-1", *req.KubernetesVersionIdentity)
			},
		},
		{
			name:    "replaces taints and keeps node labels",
			current: sampleNodePool(),
			overrides: NodePoolUpdateOverrides{
				ReplaceTaints: true,
				Taints: []kubernetes.NodeTaint{{
					Key:      "workload",
					Value:    "batch",
					Effect:   "NoSchedule",
					Operator: "Equal",
				}},
			},
			assertReq: func(t *testing.T, current *kubernetes.KubernetesNodePool, req kubernetes.UpdateKubernetesNodePool) {
				require.NotNil(t, req.NodeSettings)
				assert.Equal(t, current.NodeSettings.Labels, req.NodeSettings.Labels)
				assert.Equal(t, current.NodeSettings.Annotations, req.NodeSettings.Annotations)
				require.Len(t, req.NodeSettings.Taints, 1)
				assert.Equal(t, "workload", req.NodeSettings.Taints[0].Key)
			},
		},
		{
			name:    "replaces security groups",
			current: sampleNodePool(),
			overrides: NodePoolUpdateOverrides{
				ReplaceSecurityGroups:    true,
				SecurityGroupAttachments: []string{"sg-9"},
			},
			assertReq: func(t *testing.T, _ *kubernetes.KubernetesNodePool, req kubernetes.UpdateKubernetesNodePool) {
				assert.Equal(t, []string{"sg-9"}, req.SecurityGroupAttachments)
			},
		},
		{
			name: "requires a kubernetes version",
			current: func() *kubernetes.KubernetesNodePool {
				pool := sampleNodePool()
				pool.KubernetesVersion = nil
				return pool
			}(),
			wantErr: "kubernetes version is required",
		},
		{
			name: "requires a machine type",
			current: func() *kubernetes.KubernetesNodePool {
				pool := sampleNodePool()
				pool.MachineType.Name = ""
				return pool
			}(),
			wantErr: "machine type is required",
		},
		{
			name:    "rejects a nil node pool",
			wantErr: "node pool is required",
		},
		{
			name: "rejects a security group without an identity",
			current: func() *kubernetes.KubernetesNodePool {
				pool := sampleNodePool()
				pool.SecurityGroups = []iaas.SecurityGroup{{Name: "missing"}}
				return pool
			}(),
			wantErr: "has no identity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := BuildNodePoolUpdateRequest(tt.current, tt.overrides)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			tt.assertReq(t, tt.current, req)
		})
	}
}

func TestBuildNodePoolUpdateRequestDoesNotAliasCurrentValues(t *testing.T) {
	current := sampleNodePool()
	req, err := BuildNodePoolUpdateRequest(current, NodePoolUpdateOverrides{})
	require.NoError(t, err)

	req.Labels["team"] = "changed"
	req.NodeSettings.Labels["role"] = "changed"
	req.NodeSettings.Taints[0].Key = "changed"
	req.SecurityGroupAttachments[0] = "changed"

	assert.Equal(t, "platform", current.Labels["team"])
	assert.Equal(t, "worker", current.NodeSettings.Labels["role"])
	assert.Equal(t, "dedicated", current.NodeSettings.Taints[0].Key)
	assert.Equal(t, "sg-1", current.SecurityGroups[0].Identity)
}

func sampleNodePool() *kubernetes.KubernetesNodePool {
	return &kubernetes.KubernetesNodePool{
		Description:           "workers",
		Labels:                map[string]string{"team": "platform"},
		Annotations:           map[string]string{"note": "keep"},
		AvailabilityZone:      "nl-01a",
		EnableAutoscaling:     true,
		EnableAutoHealing:     true,
		Replicas:              2,
		MinReplicas:           1,
		MaxReplicas:           5,
		ManageNodeAllocatable: true,
		MachineType:           iaas.MachineType{Name: "pgp-medium"},
		KubernetesVersion:     &kubernetes.KubernetesVersion{Identity: "ver-1", KubernetesVersion: "1.31.1"},
		UpgradeStrategy:       kubernetes.KubernetesNodePoolUpgradeStrategyAuto,
		NodeSettings: kubernetes.KubernetesNodeSettings{
			Labels:      map[string]string{"role": "worker"},
			Annotations: map[string]string{"sched": "yes"},
			Taints: []kubernetes.NodeTaint{{
				Key:      "dedicated",
				Value:    "gpu",
				Effect:   "NoSchedule",
				Operator: "Equal",
			}},
		},
		SecurityGroups: []iaas.SecurityGroup{
			{Identity: "sg-1", Name: "one"},
			{Identity: "sg-2", Name: "two"},
		},
	}
}

func stringPtr(value string) *string { return &value }

func intPtr(value int) *int { return &value }

func boolPtr(value bool) *bool { return &value }
