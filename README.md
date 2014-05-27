gkquantile
==========

This work is based on
"Space-Efficient Online Computation of Quantile Summaries"
by M. Greenwald and S. Khanna (known as GK-algorithm)
http://infolab.stanford.edu/~datar/courses/cs361a/papers/quantiles.pdf

Also ideas for the implementation are taken from
http://www.mathcs.emory.edu/~cheung/Courses/584-StreamDB/Syllabus/08-Quantile/Greenwald.html

How to use
============

	tt := gkquantile.NewGKSummary(0.025) // where 0.025 is an accuracy. Smaller value gives a better accuracy, but will require more storage 
	tt.Add(value) // add as much values as you want
	tt.Query(0.95) - receive 0.95 quantile
	tt.Query(0.50) - receive 0.50 quantile (median)
	
	
