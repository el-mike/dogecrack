# Pitbull instance manager

## Development
1. Go to root directory `shepherd/`
2. Add .env file to the root directory:
```
SSH_USER=
SSH_PASSWORD=
SSH_DIR=

VAST_API_SECRET=

WALLET_STRING=

MONGO_INITDB_ROOT_USERNAME=
MONGO_INITDB_ROOT_PASSWORD=

```
3. Run `go install`
4. Run `docker-compose up -d`

### Fake vast.ai SSH server
Pitbull-based Docker container with running open-ssh server inside. Allows to easily test Vast.ai communication without renting actual instances.
docker-compose runs two instances of fake vast.ai server, therefore you need to specify `<number>` to connect with proper instance.

```bash
chmod +x ./tools/vast/scripts/connect_fake_vast.sh

./tools/vast/scripts/connect_fake_vast.sh <number>
# password '12345'
```

### Mongodb
Mongodb instance. Default root username/password are the ones set in `.env` file. To log in locally, run:
```bash
mongosh -u $mongoUser -p $mongoPassword
```
