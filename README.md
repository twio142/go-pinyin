# Go-Pinyin

A CLI tool for converting Chinese characters to Pinyin.

- It outputs both the original string and the Pinyin, separated with `\t`.
- Input with no Chinese characters are returned unchanged.
- Support 小鹤双拼.
- Support heteronyms (多音字).
- Full-width punctuation marks are converted to their half-width version.

## Example Use Cases

1. Used with `fzf` and `fd` to search for files with Pinyin support.

    ```bash
    fd | go-pinyin | fzf
    ```

2. Used in neovim with [`fzf-lua`](https://github.com/ibhagwan/fzf-lua)

    ```lua
    function search_file_with_pinyin()
        FzfLua.files({
            raw_cmd = 'fd ' .. FzfLua.config.setup_opts.files.fd_opts .. ' | go-pinyin',
            fzf_opts = {
                ['--delimiter'] = '\u{2002}',
                ['--with-nth'] = '{1..2}',
            },
            fn_transform = function(x)
                local file, pinyin = unpack(vim.split(x, '\t'))
                file = FzfLua.make_entry.file(file, { file_icons = true, color_icons = true }) or file
                pinyin = FzfLua.utils.ansi_codes.grey(pinyin or '')
                return string.format('%s\t%s\u{2002}%s', file, pinyin, file)
            end,
        })
    end
    ```

![screenshot](https://github.com/user-attachments/assets/6f2b0846-96d8-4986-bd3b-36bff1f3a212)

## Usage

```
$ echo '你好，世界！' | go-pinyin              # default mode
你好，世界！	ni hao , shi jie !

$ echo 'Hello, world!' | go-pinyin          # unchanged input
Hello, world!

$ echo 'Hello, 世界!' | go-pinyin -xiaohe    # xiaohe-shuangpin mode
Hello, 世界!	Hello, ui jp !

$ echo '你好，世界！' | go-pinyin -initials    # initials mode
你好，世界！	n h , s j !

$ echo '你好' | go-pinyin -only              # pinyin only
ni hao

$ echo '长期规划' | go-pinyin                 # support heteronyms
长期规划	zhang chang qi ji gui hua guo huai
```

## Installation

```bash
go install github.com/twio142/go-pinyin@latest
```

Make sure your `$GOPATH/bin` is in your `$PATH`.

## Credits

Underlying Pinyin conversion is powered by [mozillazg/go-pinyin](https://github.com/mozillazg/go-pinyin).
