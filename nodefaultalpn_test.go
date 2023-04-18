package nodefaultalpnrepro

import (
	"bytes"
	"testing"

	"github.com/coredns/coredns/plugin/test"
	coretest "github.com/coredns/coredns/test"
	"github.com/kokizzu/goproc"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestNodefaultalpn(t *testing.T) {

	// create zone1
	name, rm, err := test.TempFile(".", `$ORIGIN example.org.
@	3600 IN	SOA   sns.dns.icann.org. noc.dns.icann.org. 2017042745 7200 3600 1209600 3600

    3600 IN NS    b.iana-servers.net.

    3600 IN A 1.1.1.1
    3600 IN HTTPS 1 . no-default-alpn
`)
	if err != nil {
		t.Fatalf("Failed to create zone: %s", err)
	}
	defer rm()

	corefile := `example.org:0 {
		file ` + name + ` example.org
	}
`

	i, udp, _, err := coretest.CoreDNSServerAndPorts(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer i.Stop()
	t.Log(udp)

	t.Run(`check A`, func(t *testing.T) {
		m := new(dns.Msg)
		m.SetQuestion("example.org.", dns.TypeA)
		r, err := dns.Exchange(m, udp)
		assert.Nil(t, err)
		L.Describe(r)
		assert.Equal(t, r.Answer[0].String(), `example.org.	3600	IN	A	1.1.1.1`)
	})

	t.Run(`check HTTPS no-default-alpn`, func(t *testing.T) {
		m := new(dns.Msg)
		m.SetQuestion("example.org.", dns.TypeHTTPS)
		r, err := dns.Exchange(m, udp)
		assert.Nil(t, err)
		L.Describe(r)
		assert.Equal(t, r.Answer[0].String(), `example.org.	3600	IN	HTTPS	1 . no-default-alpn=""`)
	})

	t.Run(`check with dig`, func(t *testing.T) {
		proc := goproc.New()
		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		port := S.RightOfLast(udp, `:`)
		cmdId := proc.AddCommand(&goproc.Cmd{
			Program:    `dig`,
			Parameters: []string{`@127.0.0.1`, `-p`, port, `-t`, `https`, `example.org`},
			OnStdout: func(_ *goproc.Cmd, line string) error {
				stdout.WriteString(line)
				return nil
			},
			OnStderr: func(_ *goproc.Cmd, line string) error {
				stderr.WriteString(line)
				return nil
			},
		})
		err := proc.Start(cmdId)
		assert.Nil(t, err)
		L.Describe(stdout.String())
		L.Describe(stderr.String())
		assert.NotContains(t, stdout.String(), `bad packet: FORMERR`)
	})

}
