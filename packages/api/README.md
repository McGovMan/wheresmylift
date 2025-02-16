# WheresMyLift-API

## Services

### API

This API currently only presents one endpoint: `/v0/healthcheck`.

It relies on the following environment variables being set: WML_LOG_LEVEL, WML_HTTP_LISTEN_ADDRESS, WML_HTTP_TRUSTED_PROXY.
  - WML_LOG_LEVEL can be any of the strings named in [`config.go`](internal/config/config.go)
  - WML_HTTP_LISTEN_ADDRESS must be in the form [IP]:port, where IP is optional
  - WML_HTTP_TRUSTED_PROXY must be an IP

### Maintenance Page

The goal of the maintenance page is to run in parallel with the API on the same domain and respond with a `503 Service Temporarily Unavailable` status code and `{"message":"service unavailable"}` json responce when the API is unavailable or in an unhealthy state. This is achieved by setting the maintenance page traefik router priority to be lower than the API traefik router.

Having this maintenance page should help diagnose when the service is simply unavailable or not reachable at all.

## Testing

Run `make units` to ensure all tests pass. Run `make coverage` to ensure adaquete code coverage. `main.go` is exempt from coverage scanning and do not have any tests.

You can use `make lint` to ensure your changes conform to the code standards.

## Deploying

You'll need to have ansible installed locally. Follow these [setup instructions](https://docs.ansible.com/ansible/latest/installation_guide/installation_distros.html).

### Production Setup

The production server should be already setup with [Traefik](https://doc.traefik.io/traefik/getting-started/quick-start/). Traefik should be attached to a network called `transit-public` and have letsencrypt setup.

### API

The ansible playbook requires four variables to be set:
- env           this is the environment you wish to deploy to
- image_name    the name of the image locally or on docker hub
- endpoint      the endpoint traefik should register this container to

`ansible-playbook -e "{ env: wheresmylift.ie, image_name: stable, endpoint: api.wheresmylift.ie }" -i wheresmylift.ie, ansible/deploy.yml`

If instead you wish to deploy this locally, you can use ansible to do the same. Although you'll need to amend your `/etc/hosts` to include the domain you wish to deploy on. E.g. 

```
-- /etc/hosts --
api.localhost.wheresmylift.ie  127.0.0.1
```

`ansible-playbook -e "{ env: localhost, image_name: stable, endpoint: api.localhost.wheresmylift.ie }" -i localhost, --connection=local ansible/deploy.yml`

### Maintenance Page

The ansible playbook requires four variables to be set:
- env           this is the environment you wish to deploy to
- endpoint      the endpoint traefik should register this container to

`ansible-playbook -e "{ env: wheresmylift.ie, endpoint: api.wheresmylift.ie }" -i wheresmylift.ie, ansible/maintenance-page.yml`

If instead you wish to deploy this locally, you can use ansible to do the same. Although you'll need to amend your `/etc/hosts` to include the domain you wish to deploy on. E.g. 

```
-- /etc/hosts --
api.localhost.wheresmylift.ie  127.0.0.1
```

`ansible-playbook -e "{ env: localhost, endpoint: api.localhost.wheresmylift.ie }" -i localhost, --connection=local ansible/maintenance-page.yml`
