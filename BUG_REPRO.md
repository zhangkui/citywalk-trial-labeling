# BUG-003 Reproduction

## Issue

Combined route filters can return a list and total calculated from different conditions, while an invalid minimum-rating parameter is silently accepted.

## Steps

1. Create two routes in the same city and theme with ratings above and below the requested minimum.
2. Request the route list with city, themeId, and minRate together.
3. Compare the returned list length with total.
4. Request the route list with minRate set to a non-numeric value.

## Actual Result

- The list contains only the matching route, but total includes the lower-rated route.
- The invalid minRate request returns HTTP 200.

## Expected Result

- The list and total must use exactly the same filters.
- A non-numeric minRate must return HTTP 400.
