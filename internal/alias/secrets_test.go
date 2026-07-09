package alias

import "testing"

func TestContainsSecretDetectsLongFormFlags(t *testing.T) {
	cases := []string{
		`curl -u admin --password hunter2 https://example.com`,
		`gh auth login --token ghp_abc123`,
		`aws configure set aws_secret_access_key abc --api-key xyz`,
		`docker login --password=hunter2`,
		`curl --API-KEY=xyz https://example.com`,
	}
	for _, c := range cases {
		if !containsSecret(c) {
			t.Errorf("containsSecret(%q) = false, want true", c)
		}
	}
}

func TestContainsSecretDetectsMysqlInlinePassword(t *testing.T) {
	if !containsSecret("mysql -uroot -phunter2 mydb") {
		t.Error("expected mysql -p<password> to be detected as a secret")
	}
}

func TestContainsSecretIgnoresBarePFlagPrompt(t *testing.T) {
	// "-p" alone prompts interactively; it carries no credential.
	if containsSecret("mysql -u root -p mydb") {
		t.Error("bare -p (interactive prompt) should not be flagged")
	}
}

func TestContainsSecretIgnoresUnrelatedCommands(t *testing.T) {
	cases := []string{
		"git status --short --branch",
		"docker compose up -d",
		"kubectl get pods -n prod",
	}
	for _, c := range cases {
		if containsSecret(c) {
			t.Errorf("containsSecret(%q) = true, want false", c)
		}
	}
}

func TestContainsSecretIgnoresInlineFlagOnUnrelatedTool(t *testing.T) {
	// -p<value> is only a known inline-credential convention for the
	// allowlisted tools; on an arbitrary command it's just a short flag.
	if containsSecret("tar -pxvf archive.tar") {
		t.Error("-p on an unrelated tool should not be flagged as a secret")
	}
}

func TestContainsSecretEmptyInput(t *testing.T) {
	if containsSecret("") {
		t.Error("containsSecret(\"\") = true, want false")
	}
}
