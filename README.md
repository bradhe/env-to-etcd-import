# .env to etcd import utility

This tool is pretty straight forward: Takes a .env file (from the
[dotenv](https://github.com/bkeepers/dotenv) gem) and imports the keys in to
etcd.

## Usage

You'll need [golang](http://golang.org/doc/install) installed, of course.

``` bash
$ go install github.com/bradhe/env-to-etcd-import
$ env-to-etcd-import -node=http://127.0.0.1 -key-prefix=/my-app/my-enviro -env-file=/path/to/.env
```

That's it! To check that there is indeed data where you anticipate:

```bash
$ curl -L http://127.0.0.1:4001/api/v2/my-app/my-enviro/one-of-my-variables
{"action":"get","node":{"key":"/my-app/my-enviro/one-of-my-variables","value":"1","modifiedIndex":182,"createdIndex":20}}
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
