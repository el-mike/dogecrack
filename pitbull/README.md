# Btcrecovery-based password cracker

## Running locally

Requires NVIDIA container runtime: [Installation](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html)

```
cd pitbull/
docker build -t pitbull .
docker run --runtime=nvidia -ti pitbull
```

Once inside the container, run:
```bash
pitbull -f $FILE_URL -w $WALLET_STRING
# or
pitbull -g $GOOGLE_FILE_ID -w $WALLET_STRING
```
You can run it from anywhere in the container. Note that the file specfied in arguments will be saved in the same directory you're currently in.
