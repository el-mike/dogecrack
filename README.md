# dogecrack
Wallet password recovery toolset. Designed to automate and orchestrate password cracking processes.

- [Pitbull](https://github.com/el-mike/dogecrack/tree/master/pitbull) is a CLI tool for running and monitoring [btcrecover](https://github.com/3rdIteration/btcrecover) processes. It provides handful of useful commands. 
- [Shepherd](https://github.com/el-mike/dogecrack/tree/master/shepherd) Is a Pitbull instances orchestration service. It generates passfiles, schedules Pitbull instances using chosen GPU server provider ([vast.ai](https://vast.ai/) for example) and observe their progress and status. It comes with RPC API. 
- [Spaniel](https://github.com/el-mike/dogecrack/tree/master/spaniel) is a browser-based dashboard application, making job and instances monitoring easy.

Check tools' respective subpages for more info on how to deploy and use them.
