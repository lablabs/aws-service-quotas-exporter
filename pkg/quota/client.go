package quota

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas/types"
	"github.com/sirupsen/logrus"
)

func NewClient(log *logrus.Logger) (*Client, error) {

	cfg, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	squ := servicequotas.NewFromConfig(cfg)
	c := Client{
		log: log,
		squ: squ,
	}
	return &c, nil
}

type Client struct {
	log *logrus.Logger
	squ *servicequotas.Client
}

func (c *Client) GetQuota(ctx context.Context, serviceCode string, quotaCode string, options ...Option) (*types.ServiceQuota, error) {
	cfg := config{}
	for _, o := range options {
		o(&cfg)
	}
	if cfg.def {
		res, err := c.squ.GetAWSDefaultServiceQuota(ctx, &servicequotas.GetAWSDefaultServiceQuotaInput{
			QuotaCode:   aws.String(quotaCode),
			ServiceCode: aws.String(serviceCode),
		}, cfg.options())
		if err != nil {
			return nil, err
		}
		return res.Quota, err
	}
	res, err := c.squ.GetServiceQuota(ctx, &servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String(quotaCode),
		ServiceCode: aws.String(serviceCode),
	}, cfg.options())
	if err != nil {
		return nil, err
	}
	return res.Quota, err
}

func (c *Client) GetQuotas(ctx context.Context, serviceCode string) ([]types.ServiceQuota, error) {

	qs := make([]types.ServiceQuota, 0)
	var token *string
	for {
		res, err := c.squ.ListServiceQuotas(ctx, &servicequotas.ListServiceQuotasInput{
			ServiceCode: aws.String(serviceCode),
			NextToken:   token,
		})
		if err != nil {
			return nil, err
		}
		qs = append(qs, res.Quotas...)
		if res.NextToken == nil {
			break
		}
		token = res.NextToken
	}
	return qs, nil
}
