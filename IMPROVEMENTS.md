# Motivation
I've decided to make it OOP way for better code structuring.
Also, errors during crawling are expected, and it is not a reason to return nothing to user, so I handle it internally and return to user the result.

# Improvements
1. Add Crawler interface to follow clean architecture.
2. Add config to follow 12-Factor App approach.
3. Add logging.
4. Cache search for each rootUrl-Depth pair for specific time.
5. Consider making code concurrent. 
6. Makefile, Docker, CI/CD
7. Improve tests to cover error cases.