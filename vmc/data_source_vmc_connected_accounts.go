package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/utils"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/com/vmware/vmc/orgs/account_link/connectedAccounts"
	"gitlab.eng.vmware.com/vapi-sdk/vmc-go-sdk/vmc"
	"log"
)

func dataSourceVmcConnectedAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmcConnectedAccountsRead,

		Schema: map[string]*schema.Schema{
			"org_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Organization identifier.",
				Required:    true,
			},
			"provider_type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The cloud provider of the SDDC (AWS or ZeroCloud).",
				Optional:    true,
				Default:     "AWS",
			},
			"ids": {
				Type:        schema.TypeList,
				Description: "The corresponding connected (customer) account UUID this connection is attached to.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVmcConnectedAccountsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*vmc.Client)
	orgID := d.Get("org_id").(string)
	providerType := d.Get("provider_type").(string)
	connector, err := utils.NewVmcConnector(client.RefreshToken, "", "")
	if err != nil {
		return fmt.Errorf("Error while reading accounts from org %q: %v", orgID, err)
	}

	connectedAccountsClient := connectedAccounts.NewConnectedAccountsClientImpl(connector)
	accounts, err := connectedAccountsClient.Get(orgID, &providerType)

	ids := []string{}
	for _, account := range accounts {
		ids = append(ids, account.Id)
	}

	log.Printf("[DEBUG] Connected accounts are %v\n", accounts)

	if err != nil {
		return fmt.Errorf("Error while reading accounts from org %q: %v", orgID, err)
	}

	d.SetId(fmt.Sprintf("%s-%s", orgID, providerType))
	d.Set("ids", ids)
	return nil
}
