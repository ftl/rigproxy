# rigproxy

A proxy for the [Hamlib](https://github.com/Hamlib/Hamlib) `rigctld` server. THe proxy caches responses from all reading commands for an adjustable amount of time. This helps to reduce to load on the actual rig if there are multiple clients polling concurrently. (This actually fixes a problem with mit FT-450D if WSJT-X and CQRLog are running in concurrently.)

## Disclaimer
I develop this tool for myself and just for fun in my free time. If you find it useful, I'm happy to hear about that. If you have trouble using it, you have all the source code to fix the problem yourself (although pull requests are welcome).

## License
This tool is published under the [MIT License](https://www.tldrlegal.com/l/mit).

Copyright [Florian Thienel](http://thecodingflow.com/) 2019