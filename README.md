# Mew
## your private stock assistant
## Publish
* Mac
```
set GOOS=darwin
set GOARCH=amd64
go build
```
* Windows
```
set GOOS=windows
set GOARCH=386
go build
```
## Usage
### Auth
* You need to auth the mew to be able to let it work on your robinhood account
* We currently only supports the account with 2fa enabled
```
Y:\Projects\Mew>mew auth -u <username> -p <password> -m <2fa_code>
time="2020-05-16T23:57:10-07:00" level=info msg="welcome to use the Mew stock assistant"
time="2020-05-16T23:57:10-07:00" level=info msg="Creating config file for <username>"
```
Above will create config.yml, which stores your info to call robinhood
### Version
```
Y:\Projects\Mew>mew -v
time="2020-05-16T23:57:46-07:00" level=info msg="welcome to use the Mew stock assistant"
Mew version v0.1.2_build_2020-05-14
```
### Buy
#### Market
#### Limit
### Sell
#### Market
#### Limit
