package shellshape

import "testing"

func TestYum(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// install subcommand
		{"install single package", "yum install httpd", "yum install httpd"},
		{"install multiple packages", "yum install httpd nginx php", "yum install httpd nginx php"},
		{"install with -y", "yum install -y httpd", "yum install -y httpd"},
		{"install with assumeyes", "yum install --assumeyes httpd", "yum install --assumeyes httpd"},
		{"install quiet", "yum install -q -y vim", "yum install -q -y vim"},

		// remove/erase subcommand
		{"remove package", "yum remove httpd", "yum remove httpd"},
		{"erase package", "yum erase nginx", "yum erase nginx"},

		// update/upgrade subcommand
		{"update all", "yum update", "yum update"},
		{"update specific", "yum update httpd", "yum update httpd"},
		{"upgrade all", "yum upgrade -y", "yum upgrade -y"},

		// search subcommand
		{"search term", "yum search database", "yum search <query>"},
		{"search with flag", "yum search --all webserver", "yum search --all <query>"},

		// info subcommand
		{"info package", "yum info httpd", "yum info httpd"},

		// list subcommand
		{"list installed", "yum list installed", "yum list installed"},
		{"list available", "yum list available", "yum list available"},
		{"list all", "yum list", "yum list"},
		{"list specific", "yum list httpd", "yum list httpd"},

		// provides/whatprovides
		{"provides path", "yum provides /usr/bin/python", "yum provides <path>"},
		{"whatprovides glob", "yum whatprovides libc.so", "yum whatprovides <dotted-id>"},

		// clean subcommand
		{"clean all", "yum clean all", "yum clean all"},
		{"clean packages", "yum clean packages", "yum clean packages"},
		{"clean metadata", "yum clean metadata", "yum clean metadata"},

		// check-update
		{"check-update", "yum check-update", "yum check-update"},

		// repolist
		{"repolist", "yum repolist", "yum repolist"},
		{"repolist all", "yum repolist all", "yum repolist all"},

		// history
		{"history list", "yum history list", "yum history list"},
		{"history info id", "yum history info 5", "yum history info N"},

		// group subcommand
		{"group install", "yum group install Development-Tools", "yum group install Development-Tools"},
		{"group list", "yum group list", "yum group list"},

		// deplist
		{"deplist package", "yum deplist httpd", "yum deplist httpd"},

		// autoremove
		{"autoremove", "yum autoremove", "yum autoremove"},
		{"autoremove package", "yum autoremove httpd", "yum autoremove httpd"},

		// makecache
		{"makecache", "yum makecache", "yum makecache"},

		// flags that consume values
		{"enablerepo flag", "yum install --enablerepo=powertools gcc", "yum install --enablerepo=<val> gcc"},
		{"disablerepo flag", "yum install --disablerepo=updates httpd", "yum install --disablerepo=<val> httpd"},
		{"enablerepo separate", "yum install --enablerepo powertools gcc", "yum install --enablerepo <val> gcc"},
		{"exclude flag", "yum update --exclude kernel*", "yum update --exclude <val>"},
		{"releasever", "yum install --releasever 9 httpd", "yum install --releasever <val> httpd"},
		{"setopt", "yum install --setopt=tsflags=nodocs httpd", "yum install --setopt=<val> httpd"},
		{"config flag", "yum -c /etc/yum-custom.conf install httpd", "yum -c <path> install httpd"},
		{"installroot", "yum --installroot /mnt/sysimage install httpd", "yum --installroot <path> install httpd"},
		{"downloaddir", "yum install --downloadonly --downloaddir /tmp/rpms httpd", "yum install --downloadonly --downloaddir <path> httpd"},
		{"repo flag", "dnf install --repo=updates httpd", "dnf install --repo=<val> httpd"},

		// dnf alias
		{"dnf install", "dnf install httpd", "dnf install httpd"},
		{"dnf remove", "dnf remove nginx", "dnf remove nginx"},
		{"dnf search", "dnf search kernel", "dnf search <query>"},
		{"dnf upgrade", "dnf upgrade -y", "dnf upgrade -y"},
		{"dnf list installed", "dnf list --installed", "dnf list --installed"},

		// redirects
		{"redirect output", "yum list installed > /tmp/packages.txt", "yum list installed > <path>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different search terms collide", func(t *testing.T) {
		a := Normalize("yum search postgres")
		b := Normalize("yum search redis")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different search terms collide dnf", func(t *testing.T) {
		a := Normalize("dnf search webserver")
		b := Normalize("dnf search database")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("yum install literal-arg")
		subshell := Normalize("yum install $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})

	t.Run("subshell not collapsed dnf", func(t *testing.T) {
		benign := Normalize("dnf install literal-arg")
		subshell := Normalize("dnf install $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
