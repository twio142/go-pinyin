# Go-Pinyin

A lightweight tool for converting Chinese characters to Pinyin.

- It takes a string from standard input, and treats each line as a separate input.
- Input containing no Chinese characters are returned unchanged.
- Otherwise, it returns both the original string and the Pinyin.
- Full-width punctuation marks are converted to their half-width version.
- Supported modes:
    - Pinyin (default)
    - Pinyin initials: `-initials`
    - 小鹤双拼: `-xiaohe`
    - 小鹤双拼 initials `-xiaohe -initials`

## Example Usage

Used in `fzf` along with `fd` to search for files with Pinyin support.

```bash
fd | go-pinyin | fzf
```
