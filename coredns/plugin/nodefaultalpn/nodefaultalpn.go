package nodefaultalpn

import (
	"context"

	"github.com/coredns/caddy"
	"github.com/miekg/dns"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

var log = clog.NewWithPlugin("nodefaultalpn")

// Example is an example plugin to show how to write a plugin.
type Example struct {
	Next plugin.Handler
}

func init() {
	plugin.Register("nodefaultalpn", setup)
}

func setup(c *caddy.Controller) error {
	c.Next() // Ignore "example" and give us the next token.
	if c.NextArg() {
		// If there was another token, return an error, because we don't have any configuration.
		// Any errors returned from this setup function should be wrapped with plugin.Error, so we
		// can present a slightly nicer error message to the user.
		return plugin.Error("example", c.ArgErr())
	}

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Example{Next: next}
	})

	// All OK, return a nil error.
	return nil
}

func (e Example) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	log.Info("Received response")

	r.Answer = []dns.RR{
		&dns.HTTPS{
			SVCB: dns.SVCB{
				Hdr: dns.RR_Header{
					Name:     r.Question[0].Name,
					Rrtype:   r.Question[0].Qtype,
					Class:    r.Question[0].Qclass,
					Ttl:      1,
					Rdlength: 1,
				},
				Priority: 1,
				Target:   `.`,
				Value: []dns.SVCBKeyValue{
					&dns.SVCBNoDefaultAlpn{}, // didn't work
					//&dns.SVCBIPv4Hint{
					//	Hint: []net.IP{
					//		net.IPv4(1, 1, 1, 1),
					//	},
					//}, // works
				},
			},
		},
	}

	_ = w.WriteMsg(&dns.Msg{
		MsgHdr:   r.MsgHdr,
		Compress: false,
		Question: r.Question,
		Answer:   r.Answer,
	})

	log.Info("Override response")

	// Call next plugin (if any).
	return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
}

func (e Example) Name() string {
	return `nodefaultalpn`
}
