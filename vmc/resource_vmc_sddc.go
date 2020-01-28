/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs/sddcs"
)

func resourceSddc() *schema.Resource {
	return &schema.Resource{
		Create: resourceSddcCreate,
		Read:   resourceSddcRead,
		Update: resourceSddcUpdate,
		Delete: resourceSddcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(300 * time.Minute),
			Update: schema.DefaultTimeout(300 * time.Minute),
			Delete: schema.DefaultTimeout(180 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of this resource",
			},
			"storage_capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"sddc_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_link_sddc_config": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"customer_subnet_ids": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								// Optional: true,
							},
							Optional: true,
						},
						"connected_account_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
				ForceNew: true,
			},
			"vpc_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"num_host": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"one_node_reduced_capacity": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"sddc_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vxlan_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// TODO check the deprecation statement
			"delay_account_link": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			// TODO change default to AWS
			"provider_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "ZEROCLOUD",
			},
			"skip_creating_vxlan": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"sso_domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "vmc.local",
			},
			"sddc_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deployment_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "SingleAZ",
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "us-west-2",
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "I3_METAL",
			},
			"sddc_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vc_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nsxt_reverse_proxy_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSddcCreate(d *schema.ResourceData, m interface{}) error {
	connectorWrapper := m.(*ConnectorWrapper)
	sddcClient := orgs.NewDefaultSddcsClient(connectorWrapper)

	orgID := d.Get("org_id").(string)
	storageCapacity := d.Get("storage_capacity").(int)
	storageCapacityConverted := int64(storageCapacity)
	sddcName := d.Get("sddc_name").(string)
	vpcCidr := d.Get("vpc_cidr").(string)
	numHost := d.Get("num_host").(int)
	sddcType := d.Get("sddc_type").(string)
	oneNodeReducedCapacity := d.Get("one_node_reduced_capacity").(bool)

	if orgID == "" {
		return fmt.Errorf("org ID is a required parameter and cannot be empty")
	}
	if sddcName == "" {
		return fmt.Errorf("SDDC Name is a required parameter and cannot be empty")
	}
	if numHost == 0 {
		return fmt.Errorf("number of hosts is a required parameter and cannot be 0")
	}
	/*
		if (numHost == 1) && (oneNodeReducedCapacity == false) {
			return fmt.Errorf("one_node_reduced_capacity must be set to true with num_hosts = 1")
		}*/

	var sddcTypePtr *string
	if sddcType != "" {
		sddcTypePtr = &sddcType
	}
	vxlanSubnet := d.Get("vxlan_subnet").(string)
	delayAccountLink := d.Get("delay_account_link").(bool)
	accountLinkConfig := &model.AccountLinkConfig{
		DelayAccountLink: &delayAccountLink,
	}
	providerType := d.Get("provider_type").(string)
	skipCreatingVxlan := d.Get("skip_creating_vxlan").(bool)
	ssoDomain := d.Get("sso_domain").(string)
	sddcTemplateID := d.Get("sddc_template_id").(string)
	deploymentType := d.Get("deployment_type").(string)
	region := d.Get("region").(string)
	accountLinkSddcConfig := expandAccountLinkSddcConfig(d.Get("account_link_sddc_config").([]interface{}))
	hostInstanceType := model.HostInstanceTypes(d.Get("host_instance_type").(string))

	var awsSddcConfig = &model.AwsSddcConfig{
		StorageCapacity:        &storageCapacityConverted,
		Name:                   sddcName,
		VpcCidr:                &vpcCidr,
		NumHosts:               int64(numHost),
		OneNodeReducedCapacity: &oneNodeReducedCapacity,
		SddcType:               sddcTypePtr,
		VxlanSubnet:            &vxlanSubnet,
		AccountLinkConfig:      accountLinkConfig,
		Provider:               providerType,
		SkipCreatingVxlan:      &skipCreatingVxlan,
		AccountLinkSddcConfig:  accountLinkSddcConfig,
		SsoDomain:              &ssoDomain,
		SddcTemplateId:         &sddcTemplateID,
		DeploymentType:         &deploymentType,
		Region:                 region,
		HostInstanceType:       &hostInstanceType,
	}

	// Create a Sddc
	log.Printf("OneNodeReducedCapacity %v \n", oneNodeReducedCapacity)

	log.Printf("SDDCConfigSpew %#+v \n", *awsSddcConfig)
	task, err := sddcClient.Create(orgID, *awsSddcConfig, nil)
	if err != nil {
		return fmt.Errorf("Error while creating sddc %s: %v", sddcName, err)
	}

	// Wait until Sddc is created
	sddcID := task.ResourceId
	d.SetId(*sddcID)
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		tasksClient := orgs.NewDefaultTasksClient(connectorWrapper)
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			if err.Error() == (errors.Unauthenticated{}.Error()) {
				log.Print("Auth error", err.Error(), errors.Unauthenticated{}.Error())
				err = connectorWrapper.authenticate()
				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("Error authenticating in CSP: %s", err))
				}
				return resource.RetryableError(fmt.Errorf("Instance creation still in progress"))
			}
			return resource.NonRetryableError(fmt.Errorf("Error describing instance: %s", err))

		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("Expected instance to be created but was in state %s", *task.Status))
		}
		return resource.NonRetryableError(resourceSddcRead(d, m))
	})
	return nil
}

func resourceSddcRead(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	sddcID := d.Id()
	orgID := d.Get("org_id").(string)
	sddc, err := getSDDC(connector, orgID, sddcID)
	if err != nil {
		if err.Error() == errors.NewNotFound().Error() {
			log.Printf("SDDC with ID %s not found", sddcID)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error while getting the SDDC with ID %s,%v", sddcID, err)
	}

	if *sddc.SddcState == "DELETED" {
		log.Printf("Can't get, SDDC with ID %s is already deleted", sddc.Id)
		d.SetId("")
		return nil
	}

	d.SetId(sddc.Id)

	d.Set("name", sddc.Name)
	d.Set("updated", sddc.Updated)
	d.Set("user_id", sddc.UserId)
	d.Set("updated_by_user_id", sddc.UpdatedByUserId)
	d.Set("created", sddc.Created)
	d.Set("version", sddc.Version)
	d.Set("updated_by_user_name", sddc.UpdatedByUserName)
	d.Set("user_name", sddc.UserName)
	d.Set("org_id", sddc.OrgId)
	d.Set("sddc_type", sddc.SddcType)
	d.Set("provider", sddc.Provider)
	d.Set("account_link_state", sddc.AccountLinkState)
	d.Set("sddc_access_state", sddc.SddcAccessState)
	d.Set("sddc_type", sddc.SddcType)
	d.Set("sddc_state", sddc.SddcState)
	d.Set("one_node_reduced_capacity", sddc.OneNodeReducedCapacity)
	if sddc.ResourceConfig != nil {
		d.Set("vc_url", sddc.ResourceConfig.VcUrl)
		d.Set("cloud_username", sddc.ResourceConfig.CloudUsername)
		d.Set("cloud_password", sddc.ResourceConfig.CloudPassword)
		d.Set("nsxt_reverse_proxy_url", sddc.ResourceConfig.NsxApiPublicEndpointUrl)
	}

	return nil
}

func resourceSddcDelete(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	sddcClient := orgs.NewDefaultSddcsClient(connector)
	sddcID := d.Id()
	orgID := d.Get("org_id").(string)

	task, err := sddcClient.Delete(orgID, sddcID, nil, nil, nil)
	if err != nil {
		if err.Error() == errors.NewInvalidRequest().Error() {
			log.Printf("Can't Delete : SDDC with ID %s not found or already deleted %v", sddcID, err)
			return nil
		}
		return fmt.Errorf("Error while deleting sddc %s: %v", sddcID, err)
	}
	tasksClient := orgs.NewDefaultTasksClient(connector)
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error while deleting sddc %s: %v", sddcID, err))
		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("Expected instance to be deleted but was in state %s", *task.Status))
		}
		d.SetId("")
		return resource.NonRetryableError(nil)
	})
}

func resourceSddcUpdate(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	esxsClient := sddcs.NewDefaultEsxsClient(connector)
	sddcID := d.Id()
	orgID := d.Get("org_id").(string)

	// Add,remove hosts
	if d.HasChange("num_host") {
		oldTmp, newTmp := d.GetChange("num_host")
		oldNum := oldTmp.(int)
		newNum := newTmp.(int)

		action := "add"
		diffNum := newNum - oldNum

		if newNum < oldNum {
			action = "remove"
			diffNum = oldNum - newNum
		}

		esxConfig := model.EsxConfig{
			NumHosts: int64(diffNum),
		}

		task, err := esxsClient.Create(orgID, sddcID, esxConfig, &action)

		if err != nil {
			return fmt.Errorf("Error while updating number of host for SDDC %s: %v", sddcID, err)
		}
		tasksClient := orgs.NewDefaultTasksClient(connector)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			task, err := tasksClient.Get(orgID, task.Id)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("Error while waiting for task sddc %s: %v", task.Id, err))
			}
			if *task.Status != "FINISHED" {
				return resource.RetryableError(fmt.Errorf("Expected Host to be updated but was in state %s", *task.Status))
			}
			return resource.NonRetryableError(resourceSddcRead(d, m))
		})
		if err != nil {
			return err
		}
	}
	// Update sddc name
	if d.HasChange("sddc_name") {
		sddcClient := orgs.NewDefaultSddcsClient(connector)
		newSDDCName := d.Get("sddc_name").(string)
		sddcPatchRequest := model.SddcPatchRequest{
			Name: &newSDDCName,
		}
		sddc, err := sddcClient.Patch(orgID, sddcID, sddcPatchRequest)

		if err != nil {
			return fmt.Errorf("Error while updating sddc's name %v", err)
		}
		d.Set("sddc_name", sddc.Name)
	}
	return resourceSddcRead(d, m)
}

func expandAccountLinkSddcConfig(l []interface{}) []model.AccountLinkSddcConfig {

	if len(l) == 0 {
		return nil
	}

	var configs []model.AccountLinkSddcConfig

	for _, config := range l {
		c := config.(map[string]interface{})
		var subnetIds []string
		for _, subnetID := range c["customer_subnet_ids"].([]interface{}) {
			subnetIds = append(subnetIds, subnetID.(string))
		}
		var connectedAccId = c["connected_account_id"].(string)
		con := model.AccountLinkSddcConfig{
			CustomerSubnetIds:  subnetIds,
			ConnectedAccountId: &connectedAccId,
		}

		configs = append(configs, con)
	}
	return configs
}

func getSDDC(connector client.Connector, orgID string, sddcID string) (model.Sddc, error) {
	sddcClient := orgs.NewDefaultSddcsClient(connector)
	sddc, err := sddcClient.Get(orgID, sddcID)
	return sddc, err
}
