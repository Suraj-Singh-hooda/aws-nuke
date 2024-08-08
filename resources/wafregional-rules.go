package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafregional"
	wafregionalTypes "github.com/aws/aws-sdk-go-v2/service/wafregional/types"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WAFRegionalRule struct {
	svc     *wafregional.Client
	context context.Context

	ID   *string
	name *string
	rule *wafregionalTypes.Rule
}

func init() {
	registerV2("WAFRegionalRule", ListWAFRegionalRules)
}

func ListWAFRegionalRules(cfg *aws.Config) ([]Resource, error) {
	svc := wafregional.NewFromConfig(*cfg)
	ctx := context.TODO()
	resources := []Resource{}

	params := &wafregional.ListRulesInput{
		Limit: 50,
	}

	for {
		resp, err := svc.ListRules(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, rule := range resp.Rules {
			ruleResp, err := svc.GetRule(ctx, &wafregional.GetRuleInput{
				RuleId: rule.RuleId,
			})
			if err != nil {
				return nil, err
			}
			resources = append(resources, &WAFRegionalRule{
				svc:     svc,
				context: ctx,
				ID:      rule.RuleId,
				name:    rule.Name,
				rule:    ruleResp.Rule,
			})
		}

		if resp.NextMarker == nil {
			break
		}

		params.NextMarker = resp.NextMarker
	}

	return resources, nil
}

func (f *WAFRegionalRule) Remove() error {

	tokenOutput, err := f.svc.GetChangeToken(f.context, &wafregional.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	ruleUpdates := []wafregionalTypes.RuleUpdate{}
	for _, predicate := range f.rule.Predicates {
		ruleUpdates = append(ruleUpdates, wafregionalTypes.RuleUpdate{
			Action:    wafregionalTypes.ChangeActionDelete,
			Predicate: &predicate,
		})
	}

	if len(ruleUpdates) > 0 {
		_, err = f.svc.UpdateRule(f.context, &wafregional.UpdateRuleInput{
			ChangeToken: tokenOutput.ChangeToken,
			RuleId:      f.ID,
			Updates:     ruleUpdates,
		})

		if err != nil {
			return err
		}
	}

	tokenOutput, err = f.svc.GetChangeToken(f.context, &wafregional.GetChangeTokenInput{})
	if err != nil {
		return err
	}

	_, err = f.svc.DeleteRule(f.context, &wafregional.DeleteRuleInput{
		RuleId:      f.ID,
		ChangeToken: tokenOutput.ChangeToken,
	})

	return err
}

func (f *WAFRegionalRule) String() string {
	return *f.ID
}

func (f *WAFRegionalRule) Properties() types.Properties {
	properties := types.NewProperties()

	properties.
		Set("ID", f.ID).
		Set("Name", f.name)
	return properties
}
