package libvirt

import (
	"testing"
)

func buildTestDomain() (VirDomain, VirConnection) {
	conn := buildTestConnection()
	dom, _ := conn.LookupDomainById(1)
	return dom, conn
}

func TestGetDomainName(t *testing.T) {
	dom, conn := buildTestDomain()
	defer conn.CloseConnection()
	name, err := dom.GetName()
	if err != nil {
		t.Error(err)
		return
	}
	if name != "test" {
		t.Error("Name of active domain in test transport should be 'test'")
		return
	}
}

func TestGetDomainState(t *testing.T) {
	dom, conn := buildTestDomain()
	defer conn.CloseConnection()
	state, err := dom.GetState()
	if err != nil {
		t.Error(err)
		return
	}
	if len(state) != 2 {
		t.Error("Length of domain state should be 2")
		return
	}
	if state[0] != 1 || state[1] != 1 {
		t.Error("Domain state in test transport should be [1 1]")
		return
	}
}

func TestGetDomainUUID(t *testing.T) {
	dom, conn := buildTestDomain()
	defer conn.CloseConnection()
	_, err := dom.GetUUID()
	// how to test uuid validity?
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetDomainUUIDString(t *testing.T) {
	dom, conn := buildTestDomain()
	defer conn.CloseConnection()
	_, err := dom.GetUUIDString()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetDomainInfo(t *testing.T) {
	dom, conn := buildTestDomain()
	defer conn.CloseConnection()
	_, err := dom.GetInfo()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetDomainXMLDesc(t *testing.T) {
	dom, conn := buildTestDomain()
	defer conn.CloseConnection()
	_, err := dom.GetXMLDesc(0)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestCreateDomainSnapshotXML(t *testing.T) {
	dom, conn := buildTestDomain()
	defer conn.CloseConnection()
	_, err := dom.CreateSnapshotXML(`
		<domainsnapshot>
			<description>Test snapshot that will fail because its unsupported</description>
		</domainsnapshot>
	`, 0)
	if err == nil {
		t.Error("Snapshot should have failed due to being unsupported on test transport")
		return
	}
}

func TestSaveDomain(t *testing.T) {
	dom, conn := buildTestDomain()
	defer conn.CloseConnection()
	err := dom.Save("/tmp/libvirt-go-test.tmp")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSaveDomainFlags(t *testing.T) {
	dom, conn := buildTestDomain()
	defer conn.CloseConnection()
	err := dom.SaveFlags("/tmp/libvirt-go-test.tmp", "", 0)
	if err == nil {
		t.Error("Excected xml modification unsupported")
		return
	}
}

func TestCreateDestroyDomain(t *testing.T) {
	conn := buildTestConnection()
	defer conn.CloseConnection()
	xml := `
	<domain type="test">
		<name>test domain</name>
		<memory unit="KiB">8192</memory>
		<os>
			<type>hvm</type>
		</os>
	</domain>`
	dom, err := conn.DomainDefineXML(xml)
	if err != nil {
		t.Error(err)
		return
	}
	if err = dom.Create(); err != nil {
		t.Error(err)
		return
	}
	state, err := dom.GetState()
	if err != nil {
		t.Error(err)
		return
	}
	if state[0] != VIR_DOMAIN_RUNNING {
		t.Fatal("Domain should be running")
		return
	}
	if err = dom.Destroy(); err != nil {
		t.Error(err)
		return
	}
	state, err = dom.GetState()
	if err != nil {
		t.Error(err)
		return
	}
	if state[0] != VIR_DOMAIN_SHUTOFF {
		t.Fatal("Domain should be destroyed")
		return
	}
}

func TestShutdownDomain(t *testing.T) {
	dom, conn := buildTestDomain()
	defer conn.CloseConnection()
	if err := dom.Shutdown(); err != nil {
		t.Error(err)
		return
	}
}

func TestShutdownReboot(t *testing.T) {
	dom, conn := buildTestDomain()
	defer conn.CloseConnection()
	if err := dom.Reboot(0); err != nil {
		t.Error(err)
		return
	}
}

func TestAutostart(t *testing.T) {
	dom, conn := buildTestDomain()
	defer conn.CloseConnection()
	as, err := dom.GetAutostart()
	if err != nil {
		t.Error(err)
		return
	}
	if as {
		t.Fatal("autostart should be false")
		return
	}
	if err := dom.SetAutostart(true); err != nil {
		t.Error(err)
		return
	}
	as, err = dom.GetAutostart()
	if err != nil {
		t.Error(err)
		return
	}
	if !as {
		t.Fatal("autostart should be true")
		return
	}
}

func TestDomainIsActive(t *testing.T) {
	dom, conn := buildTestDomain()
	defer conn.CloseConnection()
	active, err := dom.IsActive()
	if err != nil {
		t.Error(err)
		return
	}
	if !active {
		t.Fatal("Domain should be active")
		return
	}
	if err := dom.Destroy(); err != nil {
		t.Error(err)
		return
	}
	active, err = dom.IsActive()
	if err != nil {
		t.Error(err)
		return
	}
	if active {
		t.Fatal("Domain should be inactive")
		return
	}
}

func TestGetMetadata(t *testing.T) {
	conn := buildTestConnection()
	defer conn.CloseConnection()
	xml := `
	<domain type="test">
		<name>test domain</name>
		<title>test domain title</title>
        <description>test domian description</description>
        <metadata>
            <appl:testapp xmlns:appl="http://testdomain/app1">testval</appl:testapp>
        </metadata>
		<memory unit="KiB">8192</memory>
		<os>
			<type>hvm</type>
		</os>
	</domain>`
	dom, err := conn.DomainDefineXML(xml)
	if err != nil {
		t.Error(err)
		return
	}
	if err = dom.Create(); err != nil {
		t.Error(err)
		return
	}

	_, err = dom.GetMetadata(VIR_DOMAIN_METADATA_DESCRIPTION, "", 0)
	if err != nil {
		t.Fatal("error in fetching domain description")
		return
	}

	_, err = dom.GetMetadata(VIR_DOMAIN_METADATA_TITLE, "", 0)
	if err != nil {
		t.Fatal("error in fetching domain title")
		return
	}

	uri := `http://testdomain/app1`
	_, err = dom.GetMetadata(VIR_DOMAIN_METADATA_ELEMENT, uri, 0)
	if err != nil {
		t.Fatal("error in fetching URI metadata")
		return
	}
}

func TestSetMetadata(t *testing.T) {
	conn := buildTestConnection()
	defer conn.CloseConnection()
	xml := `
	<domain type="test">
		<name>test domain</name>
		<title>test domain title</title>
        <description>test domian description</description>
		<memory unit="KiB">8192</memory>
		<os>
			<type>hvm</type>
		</os>
	</domain>`
	dom, err := conn.DomainDefineXML(xml)
	if err != nil {
		t.Error(err)
		return
	}
	if err = dom.Create(); err != nil {
		t.Error(err)
		return
	}

	if err := dom.SetMetadata(VIR_DOMAIN_METADATA_DESCRIPTION, "New Description", "", "", 0); err != nil {
		t.Fatal("error in setting domain description")
		return
	}

	if err := dom.SetMetadata(VIR_DOMAIN_METADATA_TITLE, "New Title", "", "", 0); err != nil {
		t.Fatal("error in setting domain title")
		return
	}

	customMetadata := `<appl:testapp xmlns:appl="http://testdomain/app1">testval</appl:testapp>`
	if err := dom.SetMetadata(VIR_DOMAIN_METADATA_ELEMENT, customMetadata, "", "", 0); err != nil {
		t.Fatal("error in fetching URI metadata")
		return
	}
}
