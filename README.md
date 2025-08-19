# Go-Pinyin

A lightweight and fast Go library for converting Chinese characters to Pinyin.

- It takes a string from standard input, and treats each line as a separate input.
- Input containing no Chinese characters are returned unchanged.
- Otherwise, it returns both the original string and the Pinyin.
- Full-width characters are converted to half-width characters.
- Supported options:
    - Pinyin (default)
    - Pinyin initials
    - 小鹤双拼
    - 小鹤双拼 initials
