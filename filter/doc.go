// Copyright 2015 Eleme Inc. All rights reserved.

/*

Package filter implements fast wildcard like filtering based on suffix tree.

Filter Suffix Tree

Build suffix tree by rules, for example:

	rules : a.b.c.d, a.b.c.e a.*.c.e
	a
	| \
	b  *
	|  |
	c  c
	|\  \
	d e  e

Rules are split into some words by '.', which wil be a node in the suffix tree.
Parent node are suffix to the word of child node. When a metric comes, the filter
search the suffix to find all matched rules at same time.

For example:
	a.b.c.e comes, the filter find 'a', then goto children node 'b' and '*', until metric
	ends or no node matched next word of metric.In the end return two rules 'a.b.c.e' and
	'a.*.c.e'.

Use suffix tree can find all matched rules at same time, it's far fast than match the rule
one by one.

Filter add and del rules by chan.

A rule in the filter can be hit by metric for the intervalHitLimit times in an interval at
most. This can help banshee limiting the metric pass.

*/
package filter
