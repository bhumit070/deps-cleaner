# Deps Cleaner

A simple CLI tool written in Go to free up disk space by deleting old dependency folders from your system. Whether it's `node_modules`, Python virtual environments, or PHP dependencies, Deps Cleaner makes cleaning up easy.

## What is it?

Deps Cleaner helps you remove unwanted dependency folders from your local system. Just give it the name of the folder, and it will take care of the rest. No more searching through your files to delete those giant, unnecessary folders. Let Deps Cleaner do the work for you.

## Why Deps Cleaner?

Let's be real—your disk is probably full of junk. You've cloned all those repos, started a bunch of projects, and then... forgot about them. It’s time to clean up. Maybe you need the space for something useful, like another repo you'll "definitely" come back to (but we both know you won’t).

## Tasks

- [ ] Add dry run flag to check which files will be deleted.
- [ ] Add file or folder to config list, from cli.

## Inspiration

I got the idea for this tool when I saw my disk slowly filling up with old projects and random repos I cloned for fun. You know, those projects you keep saying you’ll get back to? Yeah, we all have them.

I first thought of making a tool like `np-kill` just for node_modules, but then I remembered I also work with Python and PHP. So, I made this tool flexible enough to delete any dependency folder, not just JavaScript. Because why stop at cleaning up one language when you can clean up everything?

## Installation

- If you are on `mac` or `linux` run the below command.
- On windows download `deps-cleaner.exe` from [latest release](https://github.com/bhumit070/deps-cleaner/releases/latest) and you have to run the cli from that path.

```bash
curl https://raw.githubusercontent.com/bhumit070/deps-cleaner/main/install.sh | bash
```

## Usage

| Command | Description                                                                                                                                         |
| :------ | :-------------------------------------------------------------------------------------------------------------------------------------------------- |
| clean   | This is main command that will clean all folders from given folder ( by default it will use current directory but you can pass other directory too) |

- To get more info about any command just add `-h` after any command it will print info about it.
