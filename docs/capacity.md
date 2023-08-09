# Capacity

User Queries Pando to get info on a minerID

## Usage

1.Latest Raw byte capacity of a minerID

```
fce delta capacity --miner-id <miner address>
```



<img src="assets\img1.png" alt="image-20230809113652928" style="zoom: 50%;" />





2.Raw byte capacity of a minerID at a given point in time

```
fce delta capacityTime --miner-id <miner address> --unix-time <unix time>
```



<img src="assets\img2.png" alt="image-20230809113652928" style="zoom: 50%;" />



3.Amount of data sealed in a time range

```
fce delta timeRange --miner-id <miner address> --start-utime <start unix time> --end-utime <end unix time>
```



<img src="assets\img3.png" alt="image-20230809113652928" style="zoom: 50%;" />