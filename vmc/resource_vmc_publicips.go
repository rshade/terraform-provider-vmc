package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/model"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/orgs/sddcs/publicips"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/orgs/tasks"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/runtime/protocol/client"
	"reflect"
	"time"
)

func resourcePublicIP() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublicIPCreate,
		Read:   resourcePublicIPRead,
		Update: resourcePublicIPUpdate,
		Delete: resourcePublicIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Organization identifier",
			},
			"sddc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sddc Identifier",
			},
			"allocation_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP Allocation ID",
			},
			"host_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of this resource",
			},
			"private_ips": {
				Type: schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "ID of this resource",
			},
			"names": {
				Type: schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "ID of this resource",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of this resource",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of this resource",
			},
			"associated_private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of this resource",
			},
			"dnat_rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of this resource",
			},
			"snat_rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of this resource",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of action as 'attach', 'detach', 'reattach', or 'rename'",
			},
		},
	}
}

func resourcePublicIPCreate(d *schema.ResourceData, m interface{}) error {
	connector := m.(client.Connector)
	fmt.Println("Inside Create ")
	orgID := d.Get("org_id").(string)
	sddcID := d.Get("sddc_id").(string)
	hostCount := (int64)(d.Get("host_count").(int))

	var privateIPs []string
	p := reflect.ValueOf(d.Get("private_ips"))
	for i := 0; i < p.Len(); i++ {
		singleVal := p.Index(i).Elem()
		privateIPs = append(privateIPs,singleVal.String())
		fmt.Println(privateIPs)
	}

	var workloadNames []string
	s := reflect.ValueOf(d.Get("names"))
	for i := 0; i < s.Len(); i++ {
		singleVal := s.Index(i).Elem()
		workloadNames = append(workloadNames,singleVal.String())
		fmt.Println(workloadNames)
	}
	publicIPsClient := publicips.NewPublicipsClientImpl(connector)

	var sddcAllocatePublicIpSpec = &model.SddcAllocatePublicIpSpec{
		Count:      hostCount,
		PrivateIps: privateIPs,
		Names:      workloadNames,
	}

	// Create Public IP
	task, err := publicIPsClient.Create(orgID, sddcID, *sddcAllocatePublicIpSpec)
	if err != nil {
		return fmt.Errorf("error while creating public IP : %v", err)
	}

	// Wait until public IP is created
	allocationId := task.ResourceId
	fmt.Println("Inside create public IP")
	fmt.Println(*allocationId)
	fmt.Println(task.Id)
	fmt.Println(*task.ResourceType)
	fmt.Println(task.UserId)
	d.SetId(*allocationId)

	tasksClient := tasks.NewTasksClientImpl(connector)

	return resource.Retry(300*time.Minute, func() *resource.RetryError {
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error describing instance: %s", err))
		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("expected instance to be created but was in state %s", *task.Status))
		}
		return resource.NonRetryableError(resourceSddcRead(d, m))
	})
}

func resourcePublicIPRead(d *schema.ResourceData, m interface{}) error {

	publicIPClient := publicips.NewPublicipsClientImpl(m.(client.Connector))

	allocationID := d.Id()
	fmt.Println(d.Get("allocation_id").(string))
	fmt.Println("Inside PublicIP Read")
	fmt.Println(allocationID)
	orgID := d.Get("org_id").(string)
	sddcID := d.Get("sddc_id").(string)
	publicIP, err := publicIPClient.List(orgID,sddcID)
	publicIP[0].
	if err != nil {
		return fmt.Errorf("error while getting public IP details for %s: %v", allocationID, err)
	}

	d.SetId(*publicIP.AllocationId)
	d.Set("public_ip", publicIP.PublicIp)
	d.Set("name", publicIP.Name)
	d.Set("associated_private_ip", publicIP.AssociatedPrivateIp)
	d.Set("dnat_rule_id", publicIP.DnatRuleId)
	d.Set("snat_rule_id", publicIP.SnatRuleId)
	return nil
}

func resourcePublicIPDelete(d *schema.ResourceData, m interface{}) error {

	connector := m.(client.Connector)
	allocationID := d.Id()
	orgID := d.Get("org_id").(string)
	sddcID := d.Get("sddc_id").(string)
	publicIPClient := publicips.NewPublicipsClientImpl(connector)
	task, err := publicIPClient.Delete(orgID, sddcID, allocationID)
	if err != nil {
		return fmt.Errorf("error while deleting public IP %s: %v", allocationID, err)
	}

	err = WaitForTask(connector, orgID, task.Id)

	if err != nil {
		return fmt.Errorf("error while waiting for task %s: %v", task.Id, err)
	}
	d.SetId("")
	return nil
}

func resourcePublicIPUpdate(d *schema.ResourceData, m interface{}) error {
	connector := m.(client.Connector)
	publicIPClient := publicips.NewPublicipsClientImpl(connector)
	allocationID := d.Id()
	orgID := d.Get("org_id").(string)
	sddcID := d.Get("sddc_id").(string)
	action := d.Get("action").(string)
	/*
		type SddcPublicIp struct {
		        PublicIp string
		        Name *string
		        AllocationId *string
		        DnatRuleId *string
		        AssociatedPrivateIp *string
		        SnatRuleId *string
		}
	*/

	// Update public IP name
	if d.HasChange("name") {

		newPublicIPName := d.Get("name").(string)
		newSDDCPublicIP := model.SddcPublicIp{
			Name: &newPublicIPName,
		}
		_, err := publicIPClient.Update(orgID, sddcID, allocationID, action, newSDDCPublicIP)

		if err != nil {
			return fmt.Errorf("error while updating public IP's name %v", err)
		}
		d.Set("name", d.Get("name").(string))
	}
	if d.HasChange("name") {

		newPublicIPName := d.Get("name").(string)
		newSDDCPublicIP := model.SddcPublicIp{
			Name: &newPublicIPName,
		}
		_, err := publicIPClient.Update(orgID, sddcID, allocationID, action, newSDDCPublicIP)

		if err != nil {
			return fmt.Errorf("error while updating public IP's name %v", err)
		}
		d.Set("name", d.Get("name").(string))
	}

	return resourceSddcRead(d, m)
}
