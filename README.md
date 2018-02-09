# Goprsc

Goprsc is a GO client library for accessing the [Postfix Rest Server][1] V1 API. 

Goprsc currently implements all [Postfix Rest Server][1] APIs (domains, accounts, aliases, automatic sender and recipient BCC).

## Usage

```go
import "github.com/lyubenblagoev/goprsc"
```

Create a new client and use the exposed services to access different parts of the [Postfix Rest Server][1] API.

Use the DefaultClient:

```go
client := goprsc.DefaultClient
```

Create a new Client using NewClient(*http.Client):

```go
var httpClient *http.Client

...

client := goprsc.NewClient(httpClient)
```

Create a new client using NewClientWithOptions(*http.Client, ...ClientOption):

```go
client := goprsc.NewClientWithOptions(nil, HTTPSProtocolOption())
```

Client options allow changing the default protocol, host, port and user agent string using HTTPSProtocolOption(), HostOption(), PortOption() and UserAgentOption() functions. These functions return a ClientOption which changes the corresponding option in the client.

## Examples

To create a new domain:

```go
client := goprsc.DefaultClient

domainName := "example.com"

if err := client.Domains.Create(domainName); err != nil {
    fmt.Printf("Unable to create domain %s\n\n", domainName)
    return err
}
```

To get a list of all accounts in a domain:

```go

client := goprsc.DefaultClient

domainName := "example.com"

accounts, err := client.Accounts.List(domainName)
if err != nil {
    fmt.Printf("Unable to list accounts for domain %s\n\n", domainName)
    return err
}
```

Similarly you can manage other entities.

For more usage examples, you may wish to take a look at [emailctl][2] (a CLI for the [Postfix Rest Server][1]).

[1]: https://github.com/lyubenblagoev/postfix-rest-server "Postfix Rest Server"
[2]: https://github.com/lyubenblagoev/emailctl "emailctl"
