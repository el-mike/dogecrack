# Btcrecovery-based password cracker

## Local

Requires NVIDIA container runtime: [Installation](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html)

```
docker build -t pitbull .
docker run --runtime=nvidia -ti pitbull
```

## Vast.ai
1. In `Instance Configuration`, choose custom image and paste: `michalhuras/pitbull:3.0`
2. Rent a machine
3. Go to `Instances`
4. Click "connect"
5. Connect with instance via ssh

## Running btcrecovery
Once inside the container, run:
```bash
# Options
# [-w $WALLET_STRING] - wallet data extract string.
# [-f $FILE_URL] - passwordlist file that will be downloaded and tested.
# [-g $GOOGLE_FILE_ID] - Google Drive file id, can be used instead of $FILE_URL when
#   using GoogleDrive as password storage
# [-d] - runs in detached mode. You can safely close the terminal session (or log out from SSH)
#   while using it. Please note that you won't be able to track loading indicator anymore (it uses pipe buffer which is not flushed to the output with '\n'). 

#Example: 
pitbull -g $GOOGLE_FILE_ID -w $WALLET_STRING -d

# Downloads a Google Drive file with given ID, and runs btcrecover for $WALLET_STRING in detached mode.
```
You can run it from anywhere in the container. Note that the file specfied in arguments will be saved in the same directory you're currently in.

Btcrecover's output is written to `./out_btcrecover.txt` (in your current directory). Note that it does not write the loading indicator and some other "temporary" info that is usually displayed when running it attached to terminal (those informations are not "commited" to output as they do not end with newline).

If you are using `-d` flag, running `pitbull` will detach from your console as soon as passlist file is downloaded.

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
