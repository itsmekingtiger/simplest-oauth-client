# OAuth client from scratch

Simplest oauth client to test.


# How to use

```powershell
$env:CLIENT_ID="SomeWhatService"; `
$env:CLIENT_SECRET="SomeWhatSecret"; `
$env:AUTH_URL="http://auth-server/oauth/authorize"; `
$env:TOKEN_URL="http://auth-server/oauth/token"; `
go run .
```
