# Btcrecovery-based password cracker

## Local

Requires NVIDIA container runtime: [Installation](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html)

```
docker build -t pitbull .
docker run --runtime=nvidia -ti pitbull
```

## Vast.ai
1. In `Instance Configuration`, choose custom image and paste: `michalhuras/pitbull:1.0`
2. Rent a machine
3. Go to `Instances`
4. Click "connect"
5. Connect with instance via ssh

## Running btcrecovery
Once inside the container, run:
```bash
pitbull -f $FILE_URL -w $WALLET_STRING
# or
pitbull -g $GOOGLE_FILE_ID -w $WALLET_STRING
```
You can run it from anywhere in the container. Note that the file specfied in arguments will be saved in the same directory you're currently in.
