provider "vsphere" {
  user           = "administrator@vsphere.local"
  password       = "Admin!23"
  vsphere_server = "sc1-10-182-8-53.eng.vmware.com"

  # If you have a self-signed cert
  allow_unverified_ssl = true
}

data "vsphere_host" "hosts" {
  count         = "${length(var.hosts)}"
  name          = "${var.hosts[count.index]}"
  datacenter_id = "datacenter-30"
}

resource "vsphere_compute_cluster" "compute_cluster" {
  name            = "terraform-vlcm-cluster-test"
  datacenter_id   = "datacenter-30"
  host_system_ids = ["host-47",]

  drs_enabled          = false
  ha_enabled = false
  base_image_version = "7.0.0-1.0.15843807"

  # remediate
  remediate  = true
  accept_eula = true

  # export
  export_image_enabled = true
  export_software_spec = true
  export_iso_image = false
  export_offline_bundle = false

  # import
  import_image_enabled = true
  import_image_spec = "http://sc1-10-78-165-8.eng.vmware.com:9084/vum-filedownload/download?file=SOFTWARE_SPEC_1310523462.json"
}
