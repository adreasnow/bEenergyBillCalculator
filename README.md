# Bill calculator for b.energy
This module reads in a [b.energy](https://benergy.utilmate.com/) utility usage file in a `.csv` format and calculates the costs of:
- The bill to date
- The expected bill from the bill to date averaged out to 30 days
- The 7 day average exrapolated to the a 30 day average

These are soft coded with my own usage, however they can be overriden with the following flags:

```
    -er <value>   -  [0.30894]   - Electricty rate in $/KWh           
    -wr <value>   -  [18.15]     - Water rate in $/KL
    
    -gs <value>   -  [0.286]     - Gas supply charge in $/day
    -es <value>   -  [1.08661]   - Electricity supply charge in $/day
    -ws <value>   -  [0.319]     - Water supply charge in $/day
```
My building has a single daily gas usage fee, since the only gas lines that are hooked up are for the stove.