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
```
Y:\Projects\Mew>mew b -t QQQ -v 300
time="2020-05-17T00:34:35-07:00" level=info msg="welcome to use the Mew stock assistant"
time="2020-05-17T00:34:35-07:00" level=info msg="Creating rhClient..."
time="2020-05-17T00:34:35-07:00" level=info msg="Loading config..."
Please confirm the order details.
Order type: Market Buy  Security: QQQ   Quantity: 1     price: 222.74 [y/n]y
time="2020-05-17T00:34:37-07:00" level=info msg="About to place Market Buy for 1 shares of QQQ at $222.74"
time="2020-05-17T00:34:38-07:00" level=info msg="Order placed for QQQ ID dba0080c-6888-49ea-afee-c77446249681"
```
#### Limit
```
Y:\Projects\Mew>mew lb -t QQQ -l 99.5 -v 300
time="2020-05-17T00:42:54-07:00" level=info msg="welcome to use the Mew stock assistant"
time="2020-05-17T00:42:54-07:00" level=info msg="Creating rhClient..."
time="2020-05-17T00:42:54-07:00" level=info msg="Loading config..."
time="2020-05-17T00:42:56-07:00" level=info msg="quoted price is 222.740000"
time="2020-05-17T00:42:56-07:00" level=info msg="limit price is 221.630000"
Please confirm the order details.
Order type: Limit Buy   Security: QQQ   Quantity: 1     price: 221.63 [y/n]y
time="2020-05-17T00:42:58-07:00" level=info msg="About to place Limit Buy for 1 shares of QQQ at $221.63"
time="2020-05-17T00:42:58-07:00" level=info msg="Order placed with order ID 5948a01a-e491-480d-9355-bc7c1ca6a33d"
```
### Batch
```
Y:\Projects\Mew>mew lb -t @QQQ_SPY -l 99.5 -v 300
time="2020-05-17T00:43:33-07:00" level=info msg="welcome to use the Mew stock assistant"
time="2020-05-17T00:43:33-07:00" level=info msg="Creating rhClient..."
time="2020-05-17T00:43:33-07:00" level=info msg="Loading config..."
time="2020-05-17T00:43:34-07:00" level=info msg="quoted price is 222.740000"
time="2020-05-17T00:43:34-07:00" level=info msg="limit price is 221.630000"
Please confirm the order details.
Order type: Limit Buy   Security: QQQ   Quantity: 1     price: 221.63 [y/n]y
time="2020-05-17T00:43:36-07:00" level=info msg="About to place Limit Buy for 1 shares of QQQ at $221.63"
time="2020-05-17T00:43:36-07:00" level=info msg="Order placed with order ID b792c3a7-5ceb-4867-84a9-48876e2c988e"
time="2020-05-17T00:43:36-07:00" level=info msg="quoted price is 285.690000"
time="2020-05-17T00:43:36-07:00" level=info msg="limit price is 284.260000"
Please confirm the order details.
Order type: Limit Buy   Security: SPY   Quantity: 1     price: 284.26 [y/n]y
time="2020-05-17T00:43:37-07:00" level=info msg="About to place Limit Buy for 1 shares of SPY at $284.26"
time="2020-05-17T00:43:38-07:00" level=info msg="Order placed with order ID 37b4b3df-ed60-4daf-b262-ad1a5a17dbc2"
```
### Sell
#### Market
#### Limit
