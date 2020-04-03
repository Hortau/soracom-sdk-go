package soracom

import (
	"os"
	"testing"
	"time"
)

var (
	metadataClient *MetadataClient
)

var isMetaDataTestEnabled = os.Getenv("METADATA_TEST_ENABLED") != ""

func TestMetadataInit(t *testing.T) {
	if !isMetaDataTestEnabled {
		t.Skip("metadata testing is disabled")
	}

	endpoint := os.Getenv("SORACOM_METADATA_ENDPOINT")
	options := &MetadataClientOptions{
		Endpoint: endpoint,
	}
	metadataClient = NewMetadataClient(options)

	if os.Getenv("SORACOM_VERBOSE") != "" {
		metadataClient.SetVerbose(true)
	}

}

func TestMetadataGetSubscriber(t *testing.T) {
	if !isMetaDataTestEnabled {
		t.Skip("metadata testing is disabled")
	}

	_, err := metadataClient.GetSubscriber()
	if err != nil {
		t.Fatalf("GetSubscriber() failed: %v", err.Error())
	}
}

func TestMetadataUdpateSpeedClass(t *testing.T) {
	if !isMetaDataTestEnabled {
		t.Skip("metadata testing is disabled")
	}

	sub, err := metadataClient.GetSubscriber()
	if err != nil {
		t.Fatalf("GetSubscriber() failed: %v", err.Error())
	}

	originalSpeedClass := sub.SpeedClass

	sub, err = metadataClient.UpdateSpeedClass("s1.minimum")
	if err != nil {
		t.Fatalf("UpdateSpeedClass() failed: %v", err.Error())
	}
	if sub.SpeedClass != "s1.minimum" {
		t.Fatalf("Unexpected speed class")
	}

	sub, err = metadataClient.UpdateSpeedClass("s1.standard")
	if err != nil {
		t.Fatalf("UpdateSpeedClass() failed: %v", err.Error())
	}
	if sub.SpeedClass != "s1.standard" {
		t.Fatalf("Unexpected speed class")
	}

	sub, err = metadataClient.UpdateSpeedClass(originalSpeedClass)
	if err != nil {
		t.Fatalf("UpdateSpeedClass() failed: %v", err.Error())
	}
	if sub.SpeedClass != originalSpeedClass {
		t.Fatalf("Unexpected speed class")
	}
}

func TestMetadataEnableDisableTermination(t *testing.T) {
	if !isMetaDataTestEnabled {
		t.Skip("metadata testing is disabled")
	}

	sub, err := metadataClient.GetSubscriber()
	if err != nil {
		t.Fatalf("GetSubscriber() failed: %v", err.Error())
	}

	originalTerminationEnabled := sub.TerminationEnabled

	sub, err = metadataClient.EnableTermination()
	if err != nil {
		t.Fatalf("EnableTermination() failed: %v", err.Error())
	}
	if !sub.TerminationEnabled {
		t.Fatalf("Unexpected value of TerminationEnabled")
	}

	sub, err = metadataClient.DisableTermination()
	if err != nil {
		t.Fatalf("DisableTermination() failed: %v", err.Error())
	}
	if sub.TerminationEnabled {
		t.Fatalf("Unexpected value of TerminationEnabled")
	}

	if originalTerminationEnabled {
		sub, err = metadataClient.EnableTermination()
		if err != nil {
			t.Fatalf("EnableTermination() failed: %v", err.Error())
		}
	} else {
		sub, err = metadataClient.DisableTermination()
		if err != nil {
			t.Fatalf("DisableTermination() failed: %v", err.Error())
		}
	}
	if sub.TerminationEnabled != originalTerminationEnabled {
		t.Fatalf("Unexpected value of TerminationEnabled")
	}
}

func timeToUnixMilli(t time.Time) int64 {
	return t.UnixNano() / (1000 * 1000)
}

func TestMetadataSetUnsetExpiredAt(t *testing.T) {
	if !isMetaDataTestEnabled {
		t.Skip("metadata testing is disabled")
	}

	sub, err := metadataClient.GetSubscriber()
	if err != nil {
		t.Fatalf("GetSubscriber() failed: %v", err.Error())
	}

	originalExpiredAt := sub.ExpiredAt

	exp := time.Now().Add(1 * time.Hour)
	sub, err = metadataClient.SetExpiredAt(exp)
	if err != nil {
		t.Fatalf("SetExpiredAt() failed: %v", err.Error())
	}
	if sub.ExpiredAt.UnixMilli() != timeToUnixMilli(exp) {
		t.Fatalf("Unexpected ExpiredAt")
	}

	sub, err = metadataClient.UnsetExpiredAt()
	if err != nil {
		t.Fatalf("UnsetExpiredAt() failed: %v", err.Error())
	}
	if sub.ExpiredAt != nil {
		t.Fatalf("Unexpected ExpiredAt")
	}

	if originalExpiredAt != nil {
		sub, err = metadataClient.SetExpiredAt(originalExpiredAt.Time)
		if err != nil {
			t.Fatalf("SetExpiredAt() failed: %v", err.Error())
		}
		if sub.ExpiredAt.UnixMilli() != originalExpiredAt.UnixMilli() {
			t.Fatalf("Unexpected value of ExpiredAt")
		}
	}
}

func TestMetadataSetUnsetGroup(t *testing.T) {
	if !isMetaDataTestEnabled {
		t.Skip("metadata testing is disabled")
	}

	sub, err := metadataClient.GetSubscriber()
	if err != nil {
		t.Fatalf("GetSubscriber() failed: %v", err.Error())
	}

	originalGroupID := sub.GroupID
	if originalGroupID == nil {
		t.Fatalf("The calling subscriber must be in a group")
	}

	/*
		Commented out temporary because the sim being used for test
		cannot call metadata APIs after the group is unset.
			sub, err = metadataClient.UnsetGroup()
			if err != nil {
				t.Fatalf("UnsetGroup() failed %v", err.Error())
			}
			if sub.GroupID != nil {
				t.Fatalf("Unexpected GroupID")
			}
	*/

	sub, err = metadataClient.SetGroup(*originalGroupID)
	if err != nil {
		t.Fatalf("SetGroup() failed %v", err.Error())
	}
	if *sub.GroupID != *originalGroupID {
		t.Fatalf("Unexpected GroupID")
	}
}

func TestMetadataPutDeleteTags(t *testing.T) {
	if !isMetaDataTestEnabled {
		t.Skip("metadata testing is disabled")
	}

	n1 := "metadata-test-tag-name-1"
	//n1 := "metadata test tag name 1" // half width spaces must be tested
	v1 := "metadata test tag value 1"

	n2 := "メタデータテストタグ日本語"
	v2 := "metadata test tag value 2"

	sub, err := metadataClient.PutTags([]Tag{
		Tag{TagName: n1, TagValue: v1},
		Tag{TagName: n2, TagValue: v2},
	})
	if err != nil {
		t.Fatalf("PutTags() failed: %v", err.Error())
	}
	if sub.Tags[n1] != v1 {
		t.Fatalf("Unexpected tag value: %v (expected: \"%v\")", sub.Tags[n1], v1)
	}
	if sub.Tags[n2] != v2 {
		t.Fatalf("Unexpected tag value: %v (expected: \"%v\")", sub.Tags[n2], v2)
	}

	err = metadataClient.DeleteTag(n1)
	if err != nil {
		t.Fatalf("DeleteTag() failed: %v", err.Error())
	}

	sub, err = metadataClient.GetSubscriber()
	if err != nil {
		t.Fatalf("GetSubscriber() failed: %v", err.Error())
	}
	if sub.Tags[n1] != "" {
		t.Fatalf("Unexpected tag value: %v (expected: \"%v\")", sub.Tags[n1], "")
	}
	if sub.Tags[n2] != v2 {
		t.Fatalf("Unexpected tag value: %v (expected: \"%v\")", sub.Tags[n2], v2)
	}

	err = metadataClient.DeleteTag(n2)
	if err != nil {
		t.Fatalf("DeleteTag() failed: %v", err.Error())
	}

	sub, err = metadataClient.GetSubscriber()
	if err != nil {
		t.Fatalf("GetSubscriber() failed: %v", err.Error())
	}
	if sub.Tags[n1] != "" {
		t.Fatalf("Unexpected tag value: %v (expected: \"%v\")", sub.Tags[n1], "")
	}
	if sub.Tags[n2] != "" {
		t.Fatalf("Unexpected tag value: %v (expected: \"%v\")", sub.Tags[n2], "")
	}

	err = metadataClient.DeleteTag(n2)
	if err == nil {
		t.Fatalf("DeleteTag() on non-existing tag must fail")
	}
}

func TestMetadataGetUserdata(t *testing.T) {
	if !isMetaDataTestEnabled {
		t.Skip("metadata testing is disabled")
	}

	_, err := metadataClient.GetUserdata()
	if err != nil {
		t.Fatalf("GetUserdata() failed: %v", err.Error())
	}
}
