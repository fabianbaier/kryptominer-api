## API for Kryptominer.io

Official API of the Kryptominer Project.

Automated builds are available for stable releases and the unstable master branch.

## Endpoints

| Command    | Description |
|:----------:|-------------|
| GET **`/api/v1/eth/wallets/:address`** | Returns a valid json object with the current statistics of this wallet. |
| POST **`/api/v1/eth/wallets/:address`** | To update the database with the current statistics. |

### `/api/v1/eth/wallets/:address`

Response:
```
$ curl http://localhost:8787/api/v1/eth/wallets/YourAddress | jq
    {
      "activeWorkers": 3,
      "address": "YourAddress",
      "averageHashrate": 207585470,
      "balance": 200400586604772700,
      "delta": 118587103959230,
      "paid": 2e+18,
      "time": 1508877000
    }
```

### `/api/v1/eth/wallets/:address`

Inserting/Updating manually the database, e.g. for testing:

```
$ curl -H "Content-Type: application/json" -X POST -d '{ "activeWorkers": 1, "averageHashrate":210.00, "time":4, "balance":197891397387174305 }' http://localhost:8787/api/v1/eth/wallets/YourAddress
 | jq
    {
        "Status": "OK"
    }
```
## Configuration

The configuration of the API is straightforward. For more details see our webpage at [kryptominer.io](http://kryptominer.io).

### Running the database (mongodb)

Switching into the directory where the database should persist and running the containerized MongoDB Server.

```
$ docker run -d -p 27017:27017 -v $(PWD):/data/db mongo
```

Use the mongo shell to check collections:

```
$ docker exec -it <containerID> mongo
  > use Wallets
  switched to db Wallets
  > db.eth.find()
```

### Running the API (kryptominer-api)

In this root directory run:

```
$ go run main.go --ethwallet=<ETHWALLET-ADDRESS>
```

| Command    | Description |
|:----------:|-------------|
| **`--port`** | Port of the API, default: `8787`. |
| **`--db`** | URL for MongoDB Server, default: `127.0.0.1:27017`. |
| **`--ethwallet`** | ETH Address in Ethermine. |
| **`--v`** | Verbose mode, default: `false`. |

## Contribution

Thank you for considering to help out with the source code! We welcome contributions from
anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to kryptermine-api, please fork, fix, commit and send a pull request
for the maintainers to review and merge into the main code base. If you wish to submit more
complex changes though, please check up with the core devs first via our webpage [kryptominer.io](http://kryptominer.io)
to ensure those changes are in line with the general philosophy of the project and/or get some
early feedback which can make both your efforts much lighter as well as our review and merge
procedures quick and simple.

Please make sure your contributions adhere to our coding guidelines:

 * Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting) guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary) guidelines.
 * Pull requests need to be based on and opened against the `master` branch.
 * Commit messages should be prefixed with the package(s) they modify.
   * E.g. "core, ethermine: fixing worker balance rounding issue"

## License

The kryptermine-api project (i.e. all code inside this repository) under the Apache License 2.0 also
included in our repository in the `LICENSE` file.