# Vast.ai-based btcrecovery instance manager

## Fake vast.ai SSH server
Runs Pitbull-based Docker container with running open-ssh server inside. Allows to easily test Vast.ai communication without renting actual instances.

```bash
cd shepherd/
chmod +x ./scripts/connect_fake_vast.sh
docker-compose up
./scripts/connect_fake_vast.sh # password '12345'
```
