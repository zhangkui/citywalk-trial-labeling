# BUG-004 Reproduction

## Issue

Story editing can commit the story before media replacement finishes and can discard client-provided media order.

## Steps

1. Create a story with an existing media item.
2. Submit a replacement containing one valid item and one invalid oversized URL.
3. Read the story after the request fails.
4. Submit valid media with explicit order values such as 8 and 2.

## Expected Result

- Failed replacement rolls back title, content, status, and all media.
- Successful replacement preserves every submitted media order.
