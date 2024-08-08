package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafka"

	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type MSKCluster struct {
	svc     *kafka.Client
	context context.Context

	arn  string
	name string
	tags map[string]string
}

func init() {
	registerV2("MSKCluster", ListMSKCluster)
}

func ListMSKCluster(cfg *aws.Config) ([]Resource, error) {
	svc := kafka.NewFromConfig(*cfg)
	ctx := context.TODO()

	params := &kafka.ListClustersV2Input{}
	resp, err := svc.ListClustersV2(ctx, params)

	if err != nil {
		return nil, err
	}
	resources := make([]Resource, 0)
	for _, cluster := range resp.ClusterInfoList {
		resources = append(resources, &MSKCluster{
			svc:     svc,
			context: ctx,
			arn:     *cluster.ClusterArn,
			name:    *cluster.ClusterName,
			tags:    cluster.Tags,
		})

	}
	return resources, nil
}

func (m *MSKCluster) Remove() error {
	params := &kafka.DeleteClusterInput{
		ClusterArn: &m.arn,
	}

	_, err := m.svc.DeleteCluster(m.context, params)
	if err != nil {
		return err
	}
	return nil
}

func (m *MSKCluster) String() string {
	return m.arn
}

func (m *MSKCluster) Properties() types.Properties {
	properties := types.NewProperties()
	for tagKey, tagValue := range m.tags {
		properties.SetTag(aws.String(tagKey), tagValue)
	}
	properties.Set("ARN", m.arn)
	properties.Set("Name", m.name)

	return properties
}
