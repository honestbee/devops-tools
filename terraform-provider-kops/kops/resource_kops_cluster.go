package kops

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceKopsCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceKopsClusterUpdate,
		Read:   resourceKopsClusterRead,
		Update: resourceKopsClusterUpdate,
		Delete: resourceKopsClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"state_bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"manifest": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}

}

func resourceKopsClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	config := &config{
		clusterName:      d.Get("cluster_name").(string),
		kopsStatePath:    d.Get("state_bucket").(string),
		kopsManifestPath: d.Get("manifest").(string),
		sshKey:           d.Get("private_key").(string),
	}
	log.Printf("[DEBUG] Kops cluster request: %#v", config)
	out, err := updateCluster(config)
	if err != nil {
		return fmt.Errorf("Error create cluster: %s", err)
	}

	log.Printf("[DEBUG] Kops cluster response: %#v", out)
	d.SetId(config.clusterName)

	return resourceKopsClusterRead(d, meta)
}

func resourceKopsClusterDelete(d *schema.ResourceData, meta interface{}) error {
	config := &config{
		clusterName:   d.Id(),
		kopsStatePath: d.Get("state_bucket").(string),
	}
	log.Printf("[DEBUG] Delete kops cluster request: %#v", config)
	out, err := deleteCluster(config)
	if err != nil {
		return fmt.Errorf("Error delete cluster: %s", err)
	}

	log.Printf("[DEBUG] Delete kops cluster response: %#v", out)
	d.SetId(config.clusterName)

	return nil
}

func resourceKopsClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := &config{
		clusterName:   d.Id(),
		kopsStatePath: d.Get("state_bucket").(string),
	}
	log.Printf("[DEBUG] Kops cluster request: %#v", config)

	return resource.Retry(time.Duration(1)*time.Minute, func() *resource.RetryError {

		out := getCluster(config)

		if !strings.Contains(out, config.clusterName) {
			d.SetId("")
			return nil
		}

		d.Set("cluster_name", d.Id())
		return nil
	})
}
