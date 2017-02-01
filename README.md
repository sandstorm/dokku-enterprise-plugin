# dokku-enterprise-plugin
Tools and utilities for managing larger-scale Dokku deployments


## Features

### /app/nginx.conf.d

Contrary to plain dokku, you can include host nginx configuration in your deployed
app by placing the `*.conf` files inside `/app/nginx.conf.d` of your container.

The system will then extract these files and place them inside `/home/dokku/$APP/nginx.conf.d/`,
where they will be included by standard nginx configuration.

NOTE: All manually-placed files in `/home/dokku/$APP/nginx.conf.d/` will be removed.

Implementation note: This is realized using the `nginx-pre-reload` hook.

### Planned Features

#### Cloud Backup of dokku Instances

The complete instance can be exported to and imported from the cloud by extracting its essential parts:
- **Manifest**: Description of the instance environment (configuration of app, environment variables, persistent volumes, databases)
- **Application Code**: Application logic of the instance
- **Application Data**: Content of databases and persistent volumes

(todo: names to be discussed)

To cater different use cases, it will be possible to export/import only parts of an instance (e.g. only its Manifest, or Manifest + Application Code).

#### Management of dokku Instances

- Automatic backup of running instances
- Provisioning of fallback instances when primary instances fail

## Developing

### Initial Setup

```
go get github.com/DATA-DOG/godog/cmd/godog
brew install glide
glide install
```

```
./build.sh
vagrant up
```
go to `http://dokku.me` - and press "save" once.



### manual Building
```
 ./build.sh; ssh dokku@dokku.me storage:mount test /tmp:/b
 ./build.sh; ssh dokku@dokku.me
```
### Integration tests
```
./integration-test.sh
```


### Running Tests against dokku-alt
```
USE_DOKKU_ALT=1 vagrant up

cat ~/.ssh/id_rsa.pub | ssh -i .vagrant/machines/default/virtualbox/private_key vagrant@dokku.me sudo dokku access:add

ssh dokku@dokku.me help
```
# !! add your key to /root/.ssh/authorized_keys
```
./build.sh
scp -r -i .vagrant/machines/default/virtualbox/private_key bin-build/* vagrant@dokku.me:/var/lib/dokku-alt/plugins/dokku-enterprise
./integration-test.sh
```