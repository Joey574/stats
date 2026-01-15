# Stats
My personal statistics program to help be get through PHYS 121

## Download
```
go install github.com/Joey574/stats/cmd/stats@latest
```

## Usage
Stats expects a csv in the form ```x,y,truth,table```

* x -> represents our meassured value (numerical)
* y -> represents the value we want to compare against (numerical)
* truth -> represents the ground truth value (numerical)
* table -> defines the name for the table and which entries it should be groupped with (string)

### Examples
Base case
```csv
x,y,truth,table
0.4,0.2,0.3,example
0.13,0.1,0.3,example
0.32,0.5,0.3,example
```
<br>

We can skip values if they're not known or not needed (number of commas still needs to be preserved)
```csv
x,y,truth,table
0.4,,0.3,example
0.13,,0.3,example
0.32,,0.3,example
```
<br>

We can also define multiple tables
```csv
x,y,truth,table
0.4,0.2,0.3,example_1
0.13,0.1,0.3,example_2
0.32,0.5,0.3,example_2
```
<br>

Tables don't have to be defined next to each other either
```csv
x,y,truth,table
0.4,0.2,0.3,example_2
0.13,0.1,0.3,example_1
0.32,0.5,0.3,example_2
```
