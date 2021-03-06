package google

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBigtableInstance_basic(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableInstance(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableInstanceExists(
						"google_bigtable_instance.instance"),
				),
			},
		},
	})
}

func TestAccBigtableInstance_cluster(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableInstance_cluster(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableInstanceExists(
						"google_bigtable_instance.instance"),
				),
			},
		},
	})
}

func TestAccBigtableInstance_development(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableInstance_development(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableInstanceExists(
						"google_bigtable_instance.instance"),
				),
			},
		},
	})
}

func testAccCheckBigtableInstanceDestroy(s *terraform.State) error {
	var ctx = context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_bigtable_instance" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		c, err := config.bigtableClientFactory.NewInstanceAdminClient(config.Project)
		if err != nil {
			return fmt.Errorf("Error starting instance admin client. %s", err)
		}

		defer c.Close()

		_, err = c.InstanceInfo(ctx, rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Instance %s still exists.", rs.Primary.Attributes["name"])
		}
	}

	return nil
}

func testAccBigtableInstanceExists(n string) resource.TestCheckFunc {
	var ctx = context.Background()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		c, err := config.bigtableClientFactory.NewInstanceAdminClient(config.Project)
		if err != nil {
			return fmt.Errorf("Error starting instance admin client. %s", err)
		}

		defer c.Close()

		_, err = c.InstanceInfo(ctx, rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Error retrieving instance %s.", rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccBigtableInstance(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
	name = "%s"
	cluster {
		cluster_id   = "%s"
		zone         = "us-central1-b"
		num_nodes    = 3
		storage_type = "HDD"
	}
}
`, instanceName, instanceName)
}

func testAccBigtableInstance_cluster(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
	name = "%s"
	cluster {
		cluster_id   = "%s-a"
		zone         = "us-central1-b"
		num_nodes    = 3
		storage_type = "HDD"
	}
	cluster {
		cluster_id   = "%s-b"
		zone         = "us-central1-c"
		num_nodes    = 3
		storage_type = "HDD"
	}
}
`, instanceName, instanceName, instanceName)
}

func testAccBigtableInstance_development(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
	name = "%s"
	cluster {
		cluster_id    = "%s"
		zone          = "us-central1-b"
	}
	instance_type = "DEVELOPMENT"
}
`, instanceName, instanceName)
}
