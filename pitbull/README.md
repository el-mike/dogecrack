# Btcrecovery-based password cracker

## Local

Requires NVIDIA container runtime: [Installation](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html)

```
docker build -t pitbull .
docker run --runtime=nvidia -ti pitbull
```

## Vast.ai
1. In `Instance Configuration`, choose custom image and paste: `michalhuras/pitbull:2.0`
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

Btcrecover is run as a background process with output redirected to `./out_btcrecover.txt` (in your current directory). Note that it does not write the loading indicator that is usually displayed when running it attached to terminal (loading indicator are not "commited" to files as they do not end with newline).

### Status
There is a helper script:
```bash
/app/status.sh
```
It prints current status based on btcrecover output. It can returns couple of things:
* `SUCCESS: Password found: 'password'` - happens when btcrecover stops and word 'found' is found in the output. Returns the last line that should contain the found password. Exit Code: `0`
* `RUNNING: $processInfo` - happens when process is still running. Returns info from `ps`. Exit Code: `50`
* `INTERRUPTED: $lineInfo` - happens when btcrecover process is interrupted (`kill -2 $PID` for example). Contains passlist file line that was last executed. Exit Code: `51`
* `FINISHED: $additionalInfo` - any other case. Exit Code: `52`

Refer to the script itself for additional info.

### Output
You can use
```bash
/app/output.sh
```
To print current output.
