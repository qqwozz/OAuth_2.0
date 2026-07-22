#!/bin/bash

# Badge Generator Script
# Generates various badges for the README

echo "Generating badges..."

# Create assets directory if it doesn't exist
mkdir -p assets

# GitHub Stats Badge
echo "![GitHub Stats](https://github-readme-stats.vercel.app/api?username=qqwozz&show_icons=true&theme=gruvbox&bg_color=0d1117&border_color=30363d&title_color=58a6ff&text_color=c9d1d9)" > assets/badges.md

# GitHub Streak Badge
echo "![GitHub Streak](https://github-readme-streak-stats.herokuapp.com/?user=qqwozz&theme=dark&background=0d1117&stroke=30363d&ring=58a6ff&fire=58a6ff&currStreakLabel=c9d1d9&sideLabels=f0f6fc&currStreakNum=c9d1d9&sideNums=c9d1d9&dates=8b949e)" >> assets/badges.md

# Profile Trophy
echo "![Profile Trophy](https://github-profile-trophy.vercel.app/?username=qqwozz&theme=onedark&no-bg=true)" >> assets/badges.md

# Activity Graph
echo "![Activity Graph](https://github-readme-activity-graph.vercel.app/graph?username=qqwozz&bg_color=0d1117&color=58a6ff&line=30363d&point=58a6ff&area=true&area_color=58a6ff&hide_border=true)" >> assets/badges.md

echo "Badges generated in assets/badges.md"
