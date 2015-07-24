// Jensen-Shannon Divergence
//
// Fingerprints store the histogram of the window.
// Fingerprints are compared to see their JS divergence distance is less than a critical threshold.
//
// Configuration:
//  min: Minimum value of the window data.
//  max: Maximum value of the window data.
//    NOTE: The JS divergence is symmetrical and I may be able to drop the min/max requirement
//  nBins: Number of bins to use in the histogram.
//  pvalue: Standard p-value statistical threshold. Typical value is 0.05
package jsdiv
