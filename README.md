# Keyboard Annealing

Optimizes keyboard layouts for some specific metrics

- Single Finger Bigrams
- Rolls
- 3Rolls / Onehandedness
- Alternating Hands

## Install

```
https://github.com/charxii/Keyboard-Annealing.git
cd Keyboard-Annealing
go build
```

## Use

Print out stats of keyboard in `layouts.json` as well as the optimized keyboard layouts.

```
./kbannealing.exe
```

### Options

`-anneal`

Runs annealing process before printing out stats. This overwrites keyboards in `layouts.json` that starts with `000 Optimized`

`-folder`

Uses a folder that contains `monograms.txt`, `bigrams.txt`, and `trigrams.txt` for data when calculating stats and annealing. This is `CharFreqData/mt-quotes` by default.

`-text`

Uses a text file that is a word list seperated by new lines to calculate n-gram data for calculating stats and annaling. Cannot be used with -folder at the same time.

`-symbollock`

Locks symbols to QWERTY's layout in place when annaling. Might be useful if you want to stick to having symbols on the side.
