#!/bin/bash

# 获取所有标签
tags=$(git tag)

# 删除本地所有标签
git tag -d $tags

# 删除远程所有标签
for tag in $tags; do
    git push origin :refs/tags/$tag
done

echo "All tags have been deleted from both local and remote repositories."