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

#Example: 
pitbull -f $FILE_URL -w $WALLET_STRING

# Downloads a given file, and runs btcrecover for $WALLET_STRING.
```
You can run it from anywhere in the container. Note that the file specfied in arguments will be saved in the same directory you're currently in.

Pitbull runs a new terminal session with tmux, under the name "pitbull". Thanks to that, you can safely close the terminal session you started the Pitbull in (including logging out from SSH), and the process will continue to run without interruption. 
You can easily re-attach to pitbull session with:
```bash
tmux a -t "pitbull"
```
To see the live progress.

Btcrecover's output is continuously written to `progress_view.txt` file (including loading indicators). Refer to [Output](#output) to see how you can access it.

### Status
There is a helper script:
```bash
/app/status.sh
```
It prints current status based on pitbull output. It can return following statuses:
* `SUCCESS` - happens when btcrecover stops and phrase 'Password found' is found in the output. Exit Code: `0`
* `RUNNING` - happens when process is still running (i.e. pitbull process is still active). Exit Code: `50`
* `FINISHED: $additionalInfo` - any other case. Exit Code: `51`

When status script returns `SUCCESS` or `FINISHED`, check the [Output](#output) to see the btcrecover's results.

Refer to the script itself for additional info.

### Output
You can use
```bash
$pitbullDir/output.sh
# or
cat $pitbullDir/progress_view.txt
```
To get current output.
