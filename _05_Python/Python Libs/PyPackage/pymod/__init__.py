__version__ = "1.0.0"
__author__ = "Developer"

from .hello import hello

# 定义公开接口
__all__ = [
    'hello',
]

def version():
    """返回包版本"""
    return __version__