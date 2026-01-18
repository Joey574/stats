# Stats
My personal statistics program to help be get through PHYS 121

## Download
``` bash
go install github.com/Joey574/stats/cmd/stats@v0.0.6
```

## Usage

``` bash
stats -f path/to/csv/file.csv
```

You can also render to different formats, currently suports *text, color, svg, html, markdown*, defaults to text
``` bash
stats -f path/to/csv/file.csv -r svg
```

### CSV Examples
Base case
``` csv
table,label,units,x1,x2,x3
example,test 1,s,0.4,0.2,0.3
example,test 2,s,0.13,0.1,0.3
example,test 3,s,0.32,0.5,0.3
```
<br>

We can skip values if they're not known or not needed (number of commas still needs to be preserved)
``` csv
table,label,units,x1,x2,x3
example,test 1,s,0.4,0.2,
example,test 2,s,0.13,,0.3
example,test 3,s,,,0.3
```
<br>

We can also define multiple tables
``` csv
table,label,x1,y1
example_1,test 1,0.4,0.2
example_1,test 2,0.13,0.1
example_2,test 3,0.32,0.5
```
<br>

Tables don't have to be defined next to each other either
``` csv
table,label,units,x1,x2,x3
example_1,test 1,s,0.4,0.2,0.3
example_2,test 2,s,0.13,0.1,0.3
example_1,test 3,s,0.32,0.5,0.3
```
