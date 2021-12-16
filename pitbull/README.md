# Btcrecovery-based password cracker

## Running

Requires NVIDIA container runtime: [Installation](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html)

```
cd pitbull/
docker build -t pitbull .
docker run --runtime=nvidia -ti pitbull
```
Once inside the container (in `/app` directory), run:
```
./install.sh
pitbull -f $FILE_URL -w $WALLET_STRING # add -g if using GoogleDrive
```
