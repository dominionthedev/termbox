#!/usr/bin/env python3
"""
Termbox Git Monitor Service
Watches git repositories and provides status information
"""

import json
import os
import subprocess
from pathlib import Path
from typing import Dict, List, Optional, Any


class GitMonitor:
    def __init__(self, watch_dirs: List[str] = None):
        self.watch_dirs = watch_dirs or [str(Path.home())]
        
    def is_git_repo(self, path: str) -> bool:
        """Check if directory is a git repository"""
        try:
            subprocess.run(
                ["git", "-C", path, "rev-parse", "--git-dir"],
                capture_output=True,
                check=True
            )
            return True
        except subprocess.CalledProcessError:
            return False
    
    def get_repo_status(self, repo_path: str) -> Dict[str, Any]:
        """Get detailed status of a git repository"""
        if not self.is_git_repo(repo_path):
            return {}
        
        status = {
            "path": repo_path,
            "branch": self.get_current_branch(repo_path),
            "dirty": self.is_dirty(repo_path),
            "ahead": 0,
            "behind": 0,
            "untracked": 0,
            "modified": 0,
            "staged": 0
        }
        
        # Get ahead/behind counts
        try:
            result = subprocess.run(
                ["git", "-C", repo_path, "rev-list", "--left-right", "--count", "HEAD...@{upstream}"],
                capture_output=True,
                text=True,
                timeout=2
            )
            if result.returncode == 0:
                ahead, behind = map(int, result.stdout.strip().split())
                status["ahead"] = ahead
                status["behind"] = behind
        except (subprocess.TimeoutExpired, ValueError):
            pass
        
        # Get file status counts
        try:
            result = subprocess.run(
                ["git", "-C", repo_path, "status", "--porcelain"],
                capture_output=True,
                text=True,
                timeout=2
            )
            if result.returncode == 0:
                for line in result.stdout.splitlines():
                    if line.startswith("??"):
                        status["untracked"] += 1
                    elif line[0] in "MADRC":
                        status["staged"] += 1
                    elif line[1] in "MD":
                        status["modified"] += 1
        except subprocess.TimeoutExpired:
            pass
        
        return status
    
    def get_current_branch(self, repo_path: str) -> str:
        """Get current branch name"""
        try:
            result = subprocess.run(
                ["git", "-C", repo_path, "rev-parse", "--abbrev-ref", "HEAD"],
                capture_output=True,
                text=True,
                timeout=2
            )
            if result.returncode == 0:
                return result.stdout.strip()
        except subprocess.TimeoutExpired:
            pass
        return "unknown"
    
    def is_dirty(self, repo_path: str) -> bool:
        """Check if repository has uncommitted changes"""
        try:
            result = subprocess.run(
                ["git", "-C", repo_path, "status", "--porcelain"],
                capture_output=True,
                text=True,
                timeout=2
            )
            return bool(result.stdout.strip())
        except subprocess.TimeoutExpired:
            return False
    
    def find_repos(self, base_dir: str, max_depth: int = 3) -> List[str]:
        """Find all git repositories under base directory"""
        repos = []
        
        def scan_dir(path: str, depth: int = 0):
            if depth > max_depth:
                return
            
            try:
                if self.is_git_repo(path):
                    repos.append(path)
                    return  # Don't scan inside git repos
                
                for entry in os.scandir(path):
                    if entry.is_dir() and not entry.name.startswith('.'):
                        scan_dir(entry.path, depth + 1)
            except (PermissionError, OSError):
                pass
        
        scan_dir(base_dir)
        return repos
    
    def get_all_repos_status(self) -> List[Dict[str, Any]]:
        """Get status of all watched repositories"""
        all_repos = []
        
        for watch_dir in self.watch_dirs:
            repos = self.find_repos(watch_dir)
            for repo in repos:
                status = self.get_repo_status(repo)
                if status:
                    all_repos.append(status)
        
        return all_repos
    
    def export_json(self, output_path: str = "/tmp/termbox-git-status.json"):
        """Export repository status as JSON"""
        status = self.get_all_repos_status()
        
        with open(output_path, 'w') as f:
            json.dump(status, f, indent=2)
        
        return output_path


def main():
    import sys
    
    # Allow passing custom watch directories
    watch_dirs = sys.argv[1:] if len(sys.argv) > 1 else None
    
    monitor = GitMonitor(watch_dirs)
    
    # For CLI usage, print status
    if len(sys.argv) > 1 and sys.argv[1] == "--json":
        print(json.dumps(monitor.get_all_repos_status(), indent=2))
    else:
        repos = monitor.get_all_repos_status()
        for repo in repos:
            dirty_marker = "✗" if repo["dirty"] else "✓"
            ahead_behind = ""
            if repo["ahead"] > 0:
                ahead_behind += f"↑{repo['ahead']}"
            if repo["behind"] > 0:
                ahead_behind += f"↓{repo['behind']}"
            
            print(f"{dirty_marker} {repo['path']} [{repo['branch']}] {ahead_behind}")


if __name__ == "__main__":
    main()
