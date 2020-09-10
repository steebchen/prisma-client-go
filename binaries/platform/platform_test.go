package platform

import (
	"strings"
	"testing"
)

func Test_checkForExtension(t *testing.T) {
	t.Parallel()

	type args struct {
		platform string
		path     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "linux",
		args: args{
			platform: "linux",
			path:     "/some",
		},
		want: "/some",
	}, {
		name: "windows",
		args: args{
			platform: "windows",
			path:     "/some",
		},
		want: "/some.exe",
	}, {
		name: "windows with extension",
		args: args{
			platform: "windows",
			path:     "/some.gz",
		},
		want: "/some.exe.gz",
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckForExtension(tt.args.platform, tt.args.path); got != tt.want {
				t.Errorf("checkForExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseOpenSSLVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{{
		name:  "1.1",
		input: "OpenSSL 1.1.0d  1 Feb 2014",
		want:  "1.1.x",
	}, {
		name:  "1.0",
		input: "OpenSSL 1.0.2g  1 Mar 2016",
		want:  "1.0.x",
	}, {
		name:  "default to 1.1",
		input: "",
		want:  "1.1.x",
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := parseOpenSSLVersion(tt.input); got != tt.want {
				t.Errorf("parseOpenSSLVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseLinuxDistro(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{{
		name:  "default to debian",
		input: ``,
		want:  "debian",
	}, {
		name: "custom without quotes",
		input: `
			ID=fedora
		`,
		want: "rhel",
	}, {
		name: "custom with quotes",
		input: `
			ID="fedora"
		`,
		want: "rhel",
	}, {
		name: "debian",
		input: `
			PRETTY_NAME="Debian GNU/Linux 10 (buster)"
			NAME="Debian GNU/Linux"
			VERSION_ID="10"
			VERSION="10 (buster)"
			VERSION_CODENAME=buster
			ID=debian
			HOME_URL="https://www.debian.org/"
			SUPPORT_URL="https://www.debian.org/support"
			BUG_REPORT_URL="https://bugs.debian.org/"
		`,
		want: "debian",
	}, {
		name: "ubuntu",
		input: `
			NAME="Ubuntu"
			VERSION="18.04.3 LTS (Bionic Beaver)"
			ID=ubuntu
			ID_LIKE=debian
			PRETTY_NAME="Ubuntu 18.04.3 LTS"
			VERSION_ID="18.04"
			HOME_URL="https://www.ubuntu.com/"
			SUPPORT_URL="https://help.ubuntu.com/"
			BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
			PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
			VERSION_CODENAME=bionic
			UBUNTU_CODENAME=bionic
		`,
		want: "debian",
	}, {
		name: "linux mint",
		input: `
			NAME="Linux Mint"
			VERSION="18.2 (Sonya)"
			ID=linuxmint
			ID_LIKE=ubuntu
			PRETTY_NAME="Linux Mint 18.2"
			VERSION_ID="18.2"
			HOME_URL="http://www.linuxmint.com/"
			SUPPORT_URL="http://forums.linuxmint.com/"
			BUG_REPORT_URL="http://bugs.launchpad.net/linuxmint/"
			VERSION_CODENAME=sonya
			UBUNTU_CODENAME=xenial
		`,
		want: "debian",
	}, {
		name: "centos",
		input: `
			NAME="CentOS Linux"
			VERSION="8 (Core)"
			ID="centos"
			ID_LIKE="rhel fedora"
			VERSION_ID="8"
			PLATFORM_ID="platform:el8"
			PRETTY_NAME="CentOS Linux 8 (Core)"
			ANSI_COLOR="0;31"
			CPE_NAME="cpe:/o:centos:centos:8"
			HOME_URL="https://www.centos.org/"
			BUG_REPORT_URL="https://bugs.centos.org/"

			CENTOS_MANTISBT_PROJECT="CentOS-8"
			CENTOS_MANTISBT_PROJECT_VERSION="8"
			REDHAT_SUPPORT_PRODUCT="centos"
			REDHAT_SUPPORT_PRODUCT_VERSION="8"
		`,
		want: "rhel",
	}, {
		// will default to "debian"
		name: "arch",
		input: `
			NAME="Arch Linux"
			PRETTY_NAME="Arch Linux"
			ID=arch
			BUILD_ID=rolling
			ANSI_COLOR="0;36"
			HOME_URL="https://www.archlinux.org/"
			DOCUMENTATION_URL="https://wiki.archlinux.org/"
			SUPPORT_URL="https://bbs.archlinux.org/"
			BUG_REPORT_URL="https://bugs.archlinux.org/"
			LOGO=archlinux
		`,
		want: "debian",
	}, {
		name: "amazon linux 1",
		input: `
			NAME="Amazon Linux AMI"
			VERSION="2018.03"
			ID="amzn"
			ID_LIKE="rhel fedora"
			VERSION_ID="2018.03"
			PRETTY_NAME="Amazon Linux AMI 2018.03"
			ANSI_COLOR="0;33"
			CPE_NAME="cpe:/o:amazon:linux:2018.03:ga"
			HOME_URL="http://aws.amazon.com/amazon-linux-ami/"
		`,
		want: "rhel",
	}, {
		name: "amazon linux 2",
		input: `
			NAME="Amazon Linux"
			VERSION="2"
			ID="amzn"
			ID_LIKE="centos rhel fedora"
			VERSION_ID="2"
			PRETTY_NAME="Amazon Linux 2"
			ANSI_COLOR="0;33"
			CPE_NAME="cpe:2.3:o:amazon:amazon_linux:2"
			HOME_URL="https://amazonlinux.com/"
		`,
		want: "rhel",
	}, {
		name: "fedora",
		input: `
			NAME=Fedora
			VERSION="31 (Container Image)"
			ID=fedora
			VERSION_ID=31
			VERSION_CODENAME=""
			PLATFORM_ID="platform:f31"
			PRETTY_NAME="Fedora 31 (Container Image)"
			ANSI_COLOR="0;34"
			LOGO=fedora-logo-icon
			CPE_NAME="cpe:/o:fedoraproject:fedora:31"
			HOME_URL="https://fedoraproject.org/"
			DOCUMENTATION_URL="https://docs.fedoraproject.org/en-US/fedora/f31/system-administrators-guide/"
			SUPPORT_URL="https://fedoraproject.org/wiki/Communicating_and_getting_help"
			BUG_REPORT_URL="https://bugzilla.redhat.com/"
			REDHAT_BUGZILLA_PRODUCT="Fedora"
			REDHAT_BUGZILLA_PRODUCT_VERSION=31
			REDHAT_SUPPORT_PRODUCT="Fedora"
			REDHAT_SUPPORT_PRODUCT_VERSION=31
			PRIVACY_POLICY_URL="https://fedoraproject.org/wiki/Legal:PrivacyPolicy"
			VARIANT="Container Image"
			VARIANT_ID=container
		`,
		want: "rhel",
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			input := strings.ReplaceAll(tt.input, "\t", "")
			if got := parseLinuxDistro(input); got != tt.want {
				t.Errorf("parseLinuxDistro() = %v, want %v", got, tt.want)
			}
		})
	}
}
