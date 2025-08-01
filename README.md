# ndditor

Simple vim-like text editor written in Go.

## Features

- [x] insert mode
- [x] quit
- [ ] write

## Data Structure

Each line of text is stored in a gap buffer, that optimizes for the common case of inserting and deleting characters in the middle of a line.

## Installation

```bash
go get github.com/dangdungcntt/ndditor
```

## Usage

```bash
ndditor [filename]
```

### Commands

- `i`: insert mode
- `:`: command mode
- `esc`: exit to view mode

#### Command Mode Commands
- `q`: quit
- `w`: write

## License

MIT
