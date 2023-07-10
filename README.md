# aabb

I do some game development and collisions are always a topic of interest. I had a folder in my computer where I did some experiments with different approaches (brute force, bit grids, grids, hash grids, quadtrees, augmented trees, etc) for axis-aligned bounding boxes (aabb), and now I'm sharing it here.

This package is not intended to offer general nor production-level solutions, but rather:
- Help me understand the advantages, disadvantages and trade-offs between approaches by actually implementing and benchmarking them myself.
- Serve as reference material whenever I need to deal with this in the real world.

## Benchmarks

Benchmarks executed using a decent CPU (AMD Ryzen 5) on an otherwise fairly average laptop. Results sorted from best to worst speed. A few prettified results are shown on the next subsections. Raw results are provided on the last subsection.

You can also run benchmarks and tests yourself with:
```
go test -bench .
```

### Find collisions with 500 elements

| Space           | Speed |
| :-------------- | :---: |
| Grid            | 100%  |
| HashGrid        |  84%  |
| EndlessHashGrid |  66%  |
| Quadtree        |  64%  |
| BitGrid         |  48%  |
| AugmentedTree   |  33%  |
| BruteForce      |  17%  |

### Find collisions with 500 stretched elements

| Space           | Speed | Rel. to prev.    |
| :-------------- | :---: | :--------------: |
| EndlessHashGrid | 100%  |       31%        |
| Grid            |  95%  |       30%        |
| HashGrid        |  89%  |       28%        |
| BitGrid         |  84%  |       27%        |
| Quadtree        |  79%  |       25%        |
| AugmentedTree   |  53%  |       17%        |
| BruteForce      |  45%  |       14%        |

### Find collisions with 2000 elements

| Space           | Speed |
| :-------------- | :---: |
| Grid            | 100%  |
| HashGrid        |  95%  |
| Quadtree        |  75%  |
| EndlessHashGrid |  70%  |
| BitGrid         |  51%  |
| AugmentedTree   |  39%  |
| BruteForce      |  19%  |

## Raw benchmark results

| Benchmark                           |         ns/op        |
| :---------------------------------- | :------------------: |
| AugmentedTree500-12                 |      361169 ns/op    |
| AugmentedTreeStretch500-12          |      708298 ns/op    |
| AugmentedTreeHorz500-12             |      571142 ns/op    |
| AugmentedTree2000-12                |     1310849 ns/op    |
| AugmentedTreeQuarterMuts1000-12     |     1390991 ns/op    |
| AugmentedTreeStabilize2500-12       |     2547253 ns/op    |
| BitGrid500-12                       |      250445 ns/op    |
| BitGridStretch500-12                |      446176 ns/op    |
| BitGrid2000-12                      |     1011970 ns/op    |
| BitGridQuarterMuts1000-12           |     1655893 ns/op    |
| BitGridStabilize2500-12             |     1512165 ns/op    |
| BruteForce500-12                    |      692379 ns/op    |
| BruteForceStretch500-12             |      834988 ns/op    |
| BruteForceHorz500-12                |     4557665 ns/op    |
| BruteForce2000-12                   |     2695583 ns/op    |
| BruteForceQuarterMuts1000-12        |     2472993 ns/op    |
| BruteForceStabilize2500-12          |     3508488 ns/op    |
| EndlessHashGrid500-12               |      180339 ns/op    |
| EndlessHashGridStretch500-12        |      376536 ns/op    |
| EndlessHashGridHorz500-12           |      342585 ns/op    |
| EndlessHashGrid2000-12              |      738434 ns/op    |
| EndlessHashGridQuarterMuts1000-12   |      606819 ns/op    |
| EndlessHashGridStabilize2500-12     |     1282112 ns/op    |
| Grid500-12                          |      120272 ns/op    |
| GridStretch500-12                   |      394435 ns/op    |
| Grid2000-12                         |      521784 ns/op    |
| GridQuarterMuts1000-12              |      536369 ns/op    |
| GridStabilize2500-12                |      993768 ns/op    |
| HashGrid500-12                      |      142010 ns/op    |
| HashGridStretch500-12               |      419738 ns/op    |
| HashGrid2000-12                     |      549864 ns/op    |
| HashGridQuarterMuts1000-12          |      762787 ns/op    |
| HashGridStabilize2500-12            |     1013533 ns/op    |
| Quadtree500-12                      |      187935 ns/op    |
| QuadtreeStretch500-12               |      472649 ns/op    |
| QuadtreeHorz500-12                  |      551205 ns/op    |
| Quadtree2000-12                     |      694632 ns/op    |
| QuadtreeQuarterMuts1000-12          |      779952 ns/op    |
| QuadtreeStabilize2500-12            |     2140448 ns/op    |

## Trivia

**Why not use generics?**

I started this before Golang got generics, so if you were to use any of the approaches shown here, you would want to specialize the code for your own purposes. My goal wasn't to provide general implementations anyway, only to have a nice ground for studying and experimenting, so I'll continue sticking to `int` instead of `[T]`.

**Isn't the augmented tree terrible?**

It only sorts by x and uses nodes and pointers instead of slices, so... selecting it as a general solution is a very poor idea, yeah. That being said, it also has some nice properties that aren't reflected directly on the speed benchmarks. Also, although I spent a few hours trying, I couldn't optimize it like other structures, so the implementation is fairly naive and basic. Certain variants can be optimized quite a lot, but the general version is annoying to tweak. You can also do sorting for x and y with more planes for a general solution (k-d trees), but I haven't been interested in writing and testing that yet.


