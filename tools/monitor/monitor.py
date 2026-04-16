#!/usr/bin/env python3
"""
Termbox System Monitor Service
Provides real-time system stats via Unix socket
"""

import json
import os
import socket

try:
    import psutil
except ImportError:
    raise SystemExit(
        "psutil is required: pip install psutil --break-system-packages"
    )
import time
from pathlib import Path
from typing import Dict, Any


class SystemMonitor:
    def __init__(self, socket_path: str = "/tmp/termbox-monitor.sock"):
        self.socket_path = socket_path
        
    def get_system_stats(self) -> Dict[str, Any]:
        """Collect system statistics"""
        return {
            "cpu": {
                "percent": psutil.cpu_percent(interval=0.1),
                "count": psutil.cpu_count(),
                "freq": psutil.cpu_freq().current if psutil.cpu_freq() else 0
            },
            "memory": {
                "percent": psutil.virtual_memory().percent,
                "used": psutil.virtual_memory().used,
                "total": psutil.virtual_memory().total,
                "available": psutil.virtual_memory().available
            },
            "disk": {
                "percent": psutil.disk_usage('/').percent,
                "used": psutil.disk_usage('/').used,
                "total": psutil.disk_usage('/').total,
                "free": psutil.disk_usage('/').free
            },
            "network": self.get_network_stats(),
            "timestamp": time.time()
        }
    
    def get_network_stats(self) -> Dict[str, int]:
        """Get network I/O statistics"""
        net_io = psutil.net_io_counters()
        return {
            "bytes_sent": net_io.bytes_sent,
            "bytes_recv": net_io.bytes_recv,
            "packets_sent": net_io.packets_sent,
            "packets_recv": net_io.packets_recv
        }
    
    def start_server(self):
        """Start Unix socket server"""
        # Remove existing socket if present
        if os.path.exists(self.socket_path):
            os.remove(self.socket_path)
        
        # Create Unix socket
        server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        server.bind(self.socket_path)
        server.listen(5)
        
        print(f"System monitor listening on {self.socket_path}")
        
        try:
            while True:
                conn, _ = server.accept()
                try:
                    # Send current stats as JSON
                    stats = self.get_system_stats()
                    data = json.dumps(stats) + "\n"
                    conn.sendall(data.encode())
                finally:
                    conn.close()
        except KeyboardInterrupt:
            print("\nShutting down monitor...")
        finally:
            server.close()
            if os.path.exists(self.socket_path):
                os.remove(self.socket_path)


def main():
    monitor = SystemMonitor()
    monitor.start_server()


if __name__ == "__main__":
    main()
