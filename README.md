# fil-chain-extractor
A collection of services which collects Filecoin chain status data, metrics and writes them to a document database or high performance HTAP database.

## Usage

1. Initialize the fil-chain-extractor(fce) repository. Repository path would be `~/.fil-chain-extractor` if flag`--repo-path` is not set.

    ```
    fce init --repo-path <repo path>
    ```

2. Edit configs(default file path is `~/.fil-chain-extractor/config.yaml`) and setup database(mongo for the first version).

    Config example:

    ```yaml
    api:
      httpapi:
        ListenAddress: /ip4/127.0.0.1/tcp/9000
        Version: v1
      httpjsonrpcapi:
        ListenAddress: /ip4/127.0.0.1/tcp/9010
        Version: v1
    lotusinfo:
      APIAddress: /ip4/127.0.0.1/tcp/1234
      Token: ""
    storage:
      DbType: mongo
      DSN: mongodb://user:password@127.0.0.1:27017/?directConnection=true
      RW: 2
    watcher:
      MinWatchInterval: 1s
      MaxWatchInterval: 10s
    ```

3. Run daemon to listen to the chain head and extract then persist extracted data into the user-specified database.

   ```bash
   fce daemon
   ```


## Architecture

![arch](./docs/assets/arch.png)
