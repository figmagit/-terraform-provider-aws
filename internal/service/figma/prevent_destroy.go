package figma

import (
	"fmt"
	"log"
	"time"
  "os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePreventDestroy() *schema.Resource {
	return &schema.Resource{
		Create: ResourcePreventDestroyCreate,
		Read:   ResourcePreventDestroyRead,
		Update: ResourcePreventDestroyUpdate,
		Delete: ResourcePreventDestroyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"prevent_targets": {
				Type:     schema.TypeMap,
        Required: true,
			},
		},
	}
}

func ResourcePreventDestroyCreate(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func ResourcePreventDestroyRead(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func ResourcePreventDestroyUpdate(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func ResourcePreventDestroyDelete(d *schema.ResourceData, meta interface{}) error {
  if os.Getenv("TF_PREVENT_DESTROY") != "false" {
    return fmt.Errorf(
      "Destroy blocked on prevent-destroy module." +
      " We create these to act as guard rails that protect against accidental destruction of important resources." +
      " Please check your plan and make sure you are not destroying anything important." +
      " If you really mean to destroy, please set env var TF_PREVENT_DESTROY=false and re-run.")
  } else {
		log.Printf("[INFO] TF_PREVENT_DESTROY=false passed, so allowing destroy")
  }

  return nil
}
