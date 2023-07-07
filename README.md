# fil-chain-extractor
A collection of services which collects Filecoin chain status data, metrics and writes them to a document database or high performance HTAP database.

## Usage

1. Initialize the fil-chain-extractor(fce) repository. Repository path would be `~/.fil-chain-extractor` if flag`--repo-path` is not set.

    ```
    fce init --repo-path <repo path>
    ```

2. Edit configs and setup database(mongo for the first version).

3. Run daemon to listen to the chain head and extract then persist extracted data into the user-specified database.

   ```bash
   fce daemon
   ```


## Architecture

![arch](assets/arch.png)
