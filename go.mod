module github.com/aschmidt75/go-wg-wrapper

go 1.15

replace github.com/aschmidt75/go-wg-wgrapper/wgwrapper => ./pkg/wgwrapper

require golang.zx2c4.com/wireguard/wgctrl v0.0.0-20200609130330-bd2cb7843e1b
