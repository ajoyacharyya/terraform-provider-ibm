package ibm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMLbVpxHa_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMLbVpxHaConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_lb_vpx_ha.test_ha", "stay_secondary", "true"),
					testAccCheckIBMResources("ibm_lb_vpx_ha.test_ha", "primary_id",
						"ibm_lb_vpx.test_pri", "id"),
					testAccCheckIBMResources("ibm_lb_vpx_ha.test_ha", "secondary_id",
						"ibm_lb_vpx.test_sec", "id"),
				),
			},
		},
	})
}

var testAccCheckIBMLbVpxHaConfig_basic = `

resource "ibm_compute_vm_instance" "vm1" {
    hostname = "vm1"
    domain = "terraformuat.ibm.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "dal06"
    network_speed = 10
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25]
    local_disk = false
}

resource "ibm_lb_vpx" "test_pri" {
    datacenter = "dal06"
    speed = 10
    version = "10.5"
    plan = "Standard"
    ip_count = 2
    public_vlan_id = "${ibm_compute_vm_instance.vm1.public_vlan_id}"
    private_vlan_id = "${ibm_compute_vm_instance.vm1.private_vlan_id}"
    public_subnet = "${ibm_compute_vm_instance.vm1.public_subnet}"
    private_subnet = "${ibm_compute_vm_instance.vm1.private_subnet}"
}

resource "ibm_lb_vpx" "test_sec" {
    datacenter = "dal06"
    speed = 10
    version = "10.5"
    plan = "Standard"
    ip_count = 2
    public_vlan_id = "${ibm_compute_vm_instance.vm1.public_vlan_id}"
    private_vlan_id = "${ibm_compute_vm_instance.vm1.private_vlan_id}"
    public_subnet = "${ibm_compute_vm_instance.vm1.public_subnet}"
    private_subnet = "${ibm_compute_vm_instance.vm1.private_subnet}"
}

resource "ibm_lb_vpx_ha" "test_ha" {
    primary_id = "${ibm_lb_vpx.test_pri.id}"
    secondary_id = "${ibm_lb_vpx.test_sec.id}"
    stay_secondary = true
}
`