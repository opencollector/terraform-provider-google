package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func expandServiceAccountName(d TerraformResourceData, config *Config) (string, error) {
	accountId := d.Get("account_id").(string)
	serviceAccountId := d.Get("service_account_id").(string)

	if accountId != "" && serviceAccountId == "" {
		return serviceAccountFQN(accountId, d, config)
	} else if accountId == "" && serviceAccountId != "" {
		return serviceAccountId, nil
	} else {
		return "", fmt.Errorf("exactly one of account_id or service_account_id must be provided")
	}
}

func dataSourceGoogleServiceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleServiceAccountRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRFC1035Name(6, 30),
			},
			"service_account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(ServiceAccountLinkRegex),
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unique_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceAccountName, err := expandServiceAccountName(d, config)
	if err != nil {
		return err
	}

	sa, err := config.clientIAM.Projects.ServiceAccounts.Get(serviceAccountName).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Service Account %q", serviceAccountName))
	}

	d.SetId(sa.Name)
	d.Set("email", sa.Email)
	d.Set("unique_id", sa.UniqueId)
	d.Set("project", sa.ProjectId)
	d.Set("account_id", strings.Split(sa.Email, "@")[0])
	d.Set("name", sa.Name)
	d.Set("display_name", sa.DisplayName)

	return nil
}
