package notice

import "testing"

func TestGetImageKey(t *testing.T) {
	resp, err := GetImageKey("https://ipfs.io/ipfs/bafybeig6bepf5ci5fyysxlfefpjzwkfp7sarj6ed2f5a34kowgc6qenjfa/1.png", "")

	t.Logf("resp:%s", resp)
	if err != nil {
		t.Error(err)
	}
}

func TestGetTenantAccessToken(t *testing.T) {
	resp, err := GetTenantAccessToken()

	t.Logf("resp:%s", resp)

	if err != nil {
		t.Error(err)
	}
}
