
how to use this repo

```shell

cd coredns
make

./coredns -dns.port=1053

# on another terminal
dig -t https @localhost -p 1053 a whoami.example.org

;; Warning: query response not set
;; Got bad packet: FORMERR
62 bytes
d1 6d 01 20 00 01 00 01 00 00 00 01 01 61 00 00          .m...........a..
41 00 01 01 61 00 00 41 00 01 00 00 00 01 00 07          A...a..A........
00 01 00 00 02 00 00 00 00 29 04 d0 00 00 00 00          .........)......
00 0c 00 0a 00 08 ba 70 b8 64 0f c5 1b 72                .......p.d...r
;; Warning: query response not set
;; Got bad packet: FORMERR
96 bytes
71 3e 01 20 00 01 00 01 00 00 00 01 06 77 68 6f          q>...........who
61 6d 69 07 65 78 61 6d 70 6c 65 03 6f 72 67 00          ami.example.org.
00 41 00 01 06 77 68 6f 61 6d 69 07 65 78 61 6d          .A...whoami.exam
70 6c 65 03 6f 72 67 00 00 41 00 01 00 00 00 01          ple.org..A......
00 07 00 01 00 00 02 00 00 00 00 29 04 d0 00 00          ...........)....
00 00 00 0c 00 0a 00 08 ba 70 b8 64 0f c5 1b 72          .........p.d...r

```

If you change `&dns.SVCBNoDefaultAlpn{},` on `coredns/plugin/nodefaultalpn` to this code:

```
&dns.SVCBIPv4Hint{
	Hint: []net.IP{
		net.IPv4(1, 1, 1, 1),
	},
},
```

it would work:

```shell
dig -t https @localhost -p 1053 a whoami.example.org

...

;; QUESTION SECTION:
;a.                             IN      HTTPS

;; ANSWER SECTION:
a.                      1       IN      HTTPS   1 . ipv4hint=1.1.1.1
```