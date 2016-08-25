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



## Developing


./build.sh
vagrant up
go to http://dokku.me - and press "save" once.



 
 ./build.sh; ssh dokku@dokku.me storage:mount test /tmp:/b
 
 ./build.sh; ssh dokku@dokku.me 