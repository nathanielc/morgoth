// Simple fingerprinter that computes both mean and standard deviation of a window.
// Fingerprints are compared to see if the means are more than n deviations apart.
//
// Configuration:
//   The only config parameter is a 'deviations' number.
//   If the means are less than 'deviations' apart than they are considered a match.
//   Increasing 'deviations' decreases the number of anomalies detected.
package sigma
