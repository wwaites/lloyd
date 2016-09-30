Investigating the convergence properties of Lloyd's algorithm.

What if we take successive iterations and ask, what is the probability that,
an n-sided cell will become an m-sided cell on the next iteration? Use this to
construct a stochastic matrix. Is this useful?

The resulting matrix is not static, but converges about as quickly as the
distribution of polygons converges. This is measured using the six.py script:
compare the difference in probability of having a six-sided cell to the
norm of the difference in the transition probabilities.
