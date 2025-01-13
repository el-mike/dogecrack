# Pitbull instance manager

## Development
1. Go to root directory `shepherd/`
2. Add .env file to the root directory:
```
SSH_USER=
SSH_PASSWORD=
SSH_DIR=
SSH_PRIVATE_KEY=

VAST_API_SECRET=

WALLET_STRING=

MONGO_INITDB_ROOT_USERNAME=
MONGO_INITDB_ROOT_PASSWORD=

```
3. Run `go install`
4. Run `docker-compose up -d`

### Tooling
All tools should be run from `shepherd` root directory.

- Running deps services:
```bash
./tools/run_deps.sh
```

- Running deps services and apps:
```bash
./tools/run_dev.sh 
```

- Building latest `shepherd` cmd image (api + runner):
```bash
./tools/build_shepherd_image.sh 
```

- Publishing latest `shepherd` cmd image (api + runner) to Docker Hub:
```bash
./tools/publish_shepherd_image.sh 
```

### Fake vast.ai SSH server
Pitbull-based Docker container with running open-ssh server inside. Allows to easily test Vast.ai communication without renting actual instances.
`docker compose` runs two instances of fake vast.ai server, therefore you need to specify `<number>` to connect with proper instance.

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

## Setting up for Vast.ai
- Make sure your **public SSH key** is added to your Vast.ai account, in Account -> SSH Keys (so they are automatically added to newly started instances)
- add `SSH_PRIVATE_KEY` to the application's env

## Troubleshooting

### Docker images
- if building app Docker images fails due to "missing go.sum entry..." error, rebuild images with `--no-cache` flag
