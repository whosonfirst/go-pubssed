# go-httpony

Utility functions for HTTP ponies written in Go.

## Install

```
make bin
```

All the various dependencies have been included in the `vendor` directory.

## Usage

### CORS

```
import (
	"github.com/whosonfirst/go-httpony/cors"
	"net/http"
)

endpoint := "localhost:8080"
cors_enable := true
cors_allow := "*"

default_handler := func() { ... } http.Handler

// this is a standard http.HandlerFunc so assume chaining etc. here

cors_handler := cors.EnsureCORSHandler(default_handler, cors_enable, cors_allow)
http.ListenAndServe(endpoint, cors_handler)
```

### Crumb

```
import (
	"github.com/whosonfirst/go-httpony/crumb"
	"net/http"
)

// assume req is a *http.Request

ctx, _ := crumb.NewWebContext(req)

key := "G5fsBjKlsz009"
target := "admin"
length := 10
ttl := 600

c, _ := crumb.NewCrumb(ctx, key, target, length, ttl)
cr := c.Generate()

ok, err := c.Validate(cr)
```

### Crypto

```
package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-httpony/crypto"
)

func main() {

	var key = flag.String("key", "jwPsjM9rfZl73Pt0XURf0t9u8h5ZOpNT", "The key to encrypt and decrypt your text")

	flag.Parse()

	for _, text := range flag.Args() {

		c, err := crypto.NewCrypt(*key)

		if err != nil {
			panic(err)
		}

		enc, err := c.Encrypt(text)

		if err != nil {
			panic(err)
		}

		plain, err := c.Decrypt(enc)

		if err != nil {
			panic(err)
		}

		fmt.Println(text, enc, plain)
	}

}
```

### SSO

```

import (
	"github.com/whosonfirst/go-httpony/sso"
	"net/http"
)

sso_config := "/path/to/ini-config-file.cfg"
endpoint := "localhost:8080"
docroot := "www"
tls_enable := false

sso_provider, err := sso.NewSSOProvider(sso_config, endpoint, docroot, tls_enable)

if err != nil {
	panic(err)
	return
}

// this is a standard http.HandlerFunc so assume chaining etc. here

sso_handler := sso_provider.SSOHandler()
http.ListenAndServe(endpoint, sso_handler)
```

Single sign-on functionality allows your static website to act as a delegated authentication (specifically OAuth2) consumer of a different service and to use that authorization as a kind of persistent login for your own application.

When enabled a few things will happen. The first is that your web application will "grow" three new endpoints. They are:

`/signin` When visited a user will be sent to the SSO provider's OAuth2 authenticate endpoint to confirm that they want to allow your website to perform actions on their behalf.

`/auth` If a user approves your request to perform actions on their behalf they will be sent back to this endpoint and your website will complete the process to retrieve a persistent authentication token binding your website to the current user. That user's access token will be stored, encrypted, in a browser cookie whose expiration date will match the expiration date of the token itself.

`/signout` Your application can use this endpoint to "log out" a user which means that their token cookie will be removed. You will need to include a valid `crumb` parameter with the request in order for this operation to succeed. Crumbs are injected as a `data-crumb-signout` attribute in the `body` element of your application's web pages. How those values are appended to a signout URL is left for individual applications to define.

_If your web application already has URLs that map to these endpoints you will (unfortunately) need to adjust your web application accordingly. It is not currently possible to change the SSO endpoints._

On all the other HTML pages (the ones you've created for your web application) if a valid token cookie is found then it will be inserted in to page's `body` element in a `data-api-access-token` attribute. Additionally, a `data-api-endpoint` attribute (as defined in the SSO config) will be added as well as a signout "crumb". For example:

```
<body data-api-access-token="927f384c059af236a7861b87c3759ce5" data-api-endpoint="https://example.com/api/" data-crumb-signout="1467922317-42d064ad80-☃">
```

#### What now?

It is left up your web application to determine what to _do_ with these new endpoints and functionality. This includes embedding or rendering links to the `/signin` and `/signout` endpoints. Here is some sample Javascript that will inject `signin` or `signout` links when the page is loaded. Something like this mmmmmmmmmight be added to the `sso` handler in future releases...

```
<script type="text/javascript">  
window.addEventListener('load', function(e){

	var body = document.body;
	var signout_crumb = body.getAttribute("data-crumb-signout");

	if (signout_crumb){

		var signout_href = "/signout?crumb=" + encodeURIComponent(signout_crumb);

		var signout_link = document.createElement("a");
		signout_link.setAttribute("href", signout_href);
		signout_link.appendChild(document.createTextNode("sign out"));

		var signout_el = document.createElement("div");
		signout_el.setAttribute("id", "signout");
		signout_el.appendChild(signout_link);

		body.insertBefore(signout_el, body.childNodes[0]);
	}

	else {

		var signin_href = "/signin";

		var signin_link = document.createElement("a");
		signin_link.setAttribute("href", signin_href);
		signin_link.appendChild(document.createTextNode("sign in"));

		var signin_el = document.createElement("div");
		signin_el.setAttribute("id", "signin");
		signin_el.appendChild(signin_link);

		body.insertBefore(signin_el, body.childNodes[0]);
	}

});
</script>
```

The details of registering your web application, as an OAuth2 consumer, with any given third-party are outside the scope of this document. At a minimum if you are using `wof-fileserver` to run a web application locally you should make sure that the third-party service supports redirecting users to `http://localhost`

#### SSO config files

_Example:_

```
[oauth]
client_id=OAUTH2_CLIENT_ID
client_secret=OAUTH2_CLIENT_SECRET
auth_url=https://example.com/oauth2/request/
token_url=https://example.com/oauth2/token/
api_url=https://example.com/api/
scopes=write

[www]
cookie_name=sso
cookie_secret=SSO_COOKIE_SECRET
cookie_timeout=3600
crumb_secret=CRUMB_SECRET
crumb_timeout=3600
```

SSO config files are standard `ini` style config files. The `cookie_timeout` and `crumb_timeout` values are defined in seconds after which they are considered "expired". 

A `cookie_timeout` value of 0 or less will cause the SSO handler to use the expiry date of the access token (returned by the third-party service) instead. You should use this feature carefully. A `crumb_timeout` value of 0 or less will prevent the signout crumb from expiring. That's your business.

### TLS

```
import (
	"github.com/whosonfirst/go-httpony/tls"	
	"net/http"
)

// Ensures that httpony/certificates exists in your operating
// system's temporary directory and that its permissions are
// 0700. You do _not_ need to use this if you have your own
// root directory for certificates.

root, err := tls.EnsureTLSRoot()

if err != nil {
	panic(err)
}

// These are self-signed certificates so your browser _will_
// complain about them. All the usual caveats apply.

cert, key, err := tls.GenerateTLSCert(*host, root)
	
if err != nil {
	panic(err)
}

http.ListenAndServeTLS("localhost:443", cert, key, nil)
```

The details of setting up application specific HTTP handlers is left as an exercise to the reader.
