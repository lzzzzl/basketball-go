![](https://s1.ax1x.com/2022/05/26/XA4CHf.png)

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)


## Installation

You can download the latest binary from here, or you can compile from source:
```
go install github.com/lzzzzl/basketball-go@latest
```

or download the latest version pkg binaries in releases.
```
./basketball-go game -t
```

## Usage

`basketball-go` provides those main commands:
1. `game`
2. `playoff`

### Game

check **today** schedule
```
basketball-go game -t
```

check **tomorrow** schedule
```
basketball-go game -T
```

![](https://github.com/lzzzzl/GifResources/blob/main/gif/gameT.gif)

check **yesterday** shedule
```
basketball-go game -y
```

get **before** day games schedule
```
basketball-go game -b 30
```

![](https://github.com/lzzzzl/GifResources/blob/main/gif/gameB.gif)

get next day games schedule
```
basketball-go game -n 30
```

![](https://github.com/lzzzzl/GifResources/blob/main/gif/gameN.gif)

specific date to check schedule
```
basketball-go game -d 2022/5/22
```

![](https://github.com/lzzzzl/GifResources/blob/main/gif/gameD.gif)

### Playoff

check bracket of the playoff
```
basketball-go playoff -b 2022
```

![](https://github.com/lzzzzl/GifResources/blob/main/gif/playoffB.gif)

check schedule of the playoff
```
basketball-go playoff -s 2022
```

![](https://github.com/lzzzzl/GifResources/blob/main/gif/playoffS.gif)

### Team

list teams
```
basketball-go team -l
```

![](https://github.com/lzzzzl/GifResources/blob/main/gif/teamL.gif)

## License

Use of this source code is governed by the Apache 2.0 license. License that can be found in the [LICENSE](https://github.com/arsham/figurine/blob/master/LICENSE) file.

