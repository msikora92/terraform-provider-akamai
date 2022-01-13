package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourcePenaltyBox() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePenaltyBoxCreate,
		ReadContext:   resourcePenaltyBoxRead,
		UpdateContext: resourcePenaltyBoxUpdate,
		DeleteContext: resourcePenaltyBoxDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"penalty_box_protection": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"penalty_box_action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					Deny,
					Alert,
					None,
				}, false)),
			},
		},
	}
}

func resourcePenaltyBoxCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePenaltyBoxCreate")
	logger.Debugf("in resourcePenaltyBoxCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "penaltyBoxAction", m)
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	penaltyboxprotection, err := tools.GetBoolValue("penalty_box_protection", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	penaltyboxaction, err := tools.GetStringValue("penalty_box_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createPenaltyBox := appsec.UpdatePenaltyBoxRequest{
		ConfigID:             configID,
		Version:              version,
		PolicyID:             policyID,
		PenaltyBoxProtection: penaltyboxprotection,
		Action:               penaltyboxaction,
	}

	_, err = client.UpdatePenaltyBox(ctx, createPenaltyBox)
	if err != nil {
		logger.Errorf("calling 'createPenaltyBox': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", createPenaltyBox.ConfigID, createPenaltyBox.PolicyID))

	return resourcePenaltyBoxRead(ctx, d, m)
}

func resourcePenaltyBoxRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePenaltyBoxRead")
	logger.Debugf("in resourcePenaltyBoxRead")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configID, m)
	policyID := iDParts[1]

	getPenaltyBox := appsec.GetPenaltyBoxRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	penaltybox, err := client.GetPenaltyBox(ctx, getPenaltyBox)
	if err != nil {
		logger.Errorf("calling 'getPenaltyBox': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getPenaltyBox.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getPenaltyBox.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("penalty_box_protection", penaltybox.PenaltyBoxProtection); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("penalty_box_action", penaltybox.Action); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourcePenaltyBoxUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePenaltyBoxUpdate")
	logger.Debugf("in resourcePenaltyBoxUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "penaltyBoxAction", m)
	policyID := iDParts[1]
	penaltyboxprotection, err := tools.GetBoolValue("penalty_box_protection", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	penaltyboxaction, err := tools.GetStringValue("penalty_box_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updatePenaltyBox := appsec.UpdatePenaltyBoxRequest{
		ConfigID:             configID,
		Version:              version,
		PolicyID:             policyID,
		PenaltyBoxProtection: penaltyboxprotection,
		Action:               penaltyboxaction,
	}

	_, err = client.UpdatePenaltyBox(ctx, updatePenaltyBox)
	if err != nil {
		logger.Errorf("calling 'updatePenaltyBox': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourcePenaltyBoxRead(ctx, d, m)
}

func resourcePenaltyBoxDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePenaltyBoxDelete")
	logger.Debugf("in resourcePenaltyBoxDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "penaltyBoxAction", m)
	policyID := iDParts[1]

	removePenaltyBox := appsec.UpdatePenaltyBoxRequest{
		ConfigID:             configID,
		Version:              version,
		PolicyID:             policyID,
		PenaltyBoxProtection: false,
		Action:               "none",
	}

	_, err = client.UpdatePenaltyBox(ctx, removePenaltyBox)
	if err != nil {
		logger.Errorf("calling 'removePenaltyBox': %s", err.Error())
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}
