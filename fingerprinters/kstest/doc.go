// Kolmogorov–Smirnov test.
// https://en.wikipedia.org/wiki/Kolmogorov%E2%80%93Smirnov_test
//
// The fingerprint is the cummulative distribution of the window.
// The fingerprints are compared by computing the largest distance between the cummulative distribution functions and comparing to a critical value.
//
// Configuration:
//  The only parameter is a confidence level.
//  Valid values are from 0-5.
//  The level maps to a list of predefined critical values for the KS test.
//  Increasing 'confidence' decreases the number of anomalies detected.
//
package kstest
